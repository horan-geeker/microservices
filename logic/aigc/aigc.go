package aigc

import (
	"context"
	"fmt"
	"microservices/cache"
	"microservices/entity/response"
	"microservices/logic/auth"
	"microservices/pkg/log"

	entity "microservices/entity/model"
	"microservices/model"
	"microservices/service"

	"microservices/entity/consts"

	"google.golang.org/genai"
	"time"
)

type GenerationLogic interface {
	Generate(ctx context.Context, apiKey, model, prompt string) (uint64, error)
	List(ctx context.Context, uid int, page, size int) (*response.GetGenerationsResp, error)
	Detail(ctx context.Context, uid int, id uint64) (*response.GenerationDetailResp, error)
	GetResult(ctx context.Context, uid int, id uint64) (*response.GenerationResponse, error)
}

type generationLogic struct {
	model     model.Factory
	cache     cache.Factory
	service   service.Factory
	authLogic auth.Logic
}

func NewGenerationLogic(model model.Factory, cache cache.Factory, service service.Factory) GenerationLogic {
	return &generationLogic{
		model:     model,
		cache:     cache,
		service:   service,
		authLogic: auth.NewAuth(model, cache, service),
	}
}

func (l *generationLogic) Generate(ctx context.Context, apiKey, model, prompt string) (uint64, error) {
	authUser, err := l.authLogic.GetAuthUser(ctx)
	if err != nil {
		return 0, err
	}
	// 0. Check Lock
	locked, err := l.cache.AIGC().SetGenerationLock(ctx, authUser.Uid)
	if err != nil {
		return 0, err
	}
	if !locked {
		return 0, fmt.Errorf("Generate already in progress")
	}
	// Create Generation Record
	gen := &entity.Generation{
		UserID:    authUser.Uid,
		Type:      "text",
		ModelName: model,
		Prompt:    prompt, // todo 补全用户 prompt
		Status:    entity.GenerationStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := l.model.Generation().Create(ctx, gen); err != nil {
		return 0, err
	}

	// Async call
	go func(genID uint64, prompt string, uid int) {
		// Ensure lock is released even if panic
		defer func() {
			_ = l.cache.AIGC().ReleaseGenerationLock(context.Background(), uid)
		}()
		resp, err := l.service.Google().GenerateContent(
			ctx,
			apiKey,
			model,
			[]*genai.Part{
				{Text: prompt},
			},
			&genai.Content{
				Parts: []*genai.Part{
					{Text: consts.SYSTEM_INSTRUCTION},
				},
			},
		)

		status := entity.GenerationStatusFailed
		content := ""
		if err == nil && resp != nil && len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			status = entity.GenerationStatusSuccess
			content = resp.Candidates[0].Content.Parts[0].Text
		}

		if err := l.model.Generation().Update(ctx, genID, map[string]interface{}{
			"status":       status,
			"content_text": content,
			"updated_at":   time.Now(),
		}); err != nil {
			log.Error(ctx, "sql-error", err, nil)
		}
	}(gen.ID, prompt, authUser.Uid)

	return gen.ID, nil
}

func (l *generationLogic) List(ctx context.Context, uid int, page, size int) (*response.GetGenerationsResp, error) {
	gens, total, err := l.model.Generation().ListByUserID(ctx, uid, page, size)
	if err != nil {
		return nil, err
	}

	var genIDs []uint64
	for _, g := range gens {
		genIDs = append(genIDs, g.ID)
	}

	filesMap, err := l.model.Generation().GetFilesByGenerationIDs(ctx, genIDs)
	if err != nil {
		// Log error but proceed? Or fail. Let's proceed with empty files.
		log.Error(ctx, "get-generation-files-error", err, nil)
	}

	list := make([]*response.GenerationSummary, 0, len(gens))
	for _, g := range gens {
		summary := &response.GenerationSummary{
			ID:        g.ID,
			Status:    g.Status,
			CreatedAt: g.CreatedAt,
		}

		if files, ok := filesMap[g.ID]; ok {
			// Type 1: Input (Original), Type 2: Output (Generated)
			if f, ok := files[entity.GenerationFileTypeInput]; ok {
				// Sign URL
				signedUrl, err := l.service.Cloudflare().SignUrl(ctx, f.FileKey, 3600*24)
				if err == nil {
					f.Url = signedUrl
				}
				summary.OriginalFile = f
			}
			if f, ok := files[entity.GenerationFileTypeOutput]; ok {
				// Sign URL
				signedUrl, err := l.service.Cloudflare().SignUrl(ctx, f.FileKey, 3600*24)
				if err == nil {
					f.Url = signedUrl
				}
				summary.GeneratedFile = f
			}
		}
		list = append(list, summary)
	}

	return &response.GetGenerationsResp{
		List:  list,
		Total: total,
	}, nil
}

func (l *generationLogic) Detail(ctx context.Context, uid int, id uint64) (*response.GenerationDetailResp, error) {
	gen, err := l.model.Generation().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if gen.UserID != uid {
		return nil, fmt.Errorf("permission denied")
	}

	return &response.GenerationDetailResp{
		ID:          gen.ID,
		Status:      gen.Status,
		ContentText: gen.ContentText,
		CreatedAt:   gen.CreatedAt,
		UpdatedAt:   gen.UpdatedAt,
	}, nil
}

func (l *generationLogic) GetResult(ctx context.Context, uid int, id uint64) (*response.GenerationResponse, error) {
	gen, err := l.model.Generation().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if gen.UserID != uid {
		return nil, fmt.Errorf("permission denied")
	}

	var imageUrl string
	if gen.Status == entity.GenerationStatusSuccess {
		// Get Output File
		file, err := l.model.Generation().GetOutputFileByGenID(ctx, id)
		if err == nil && file != nil {
			// Sign URL
			signedUrl, err := l.service.Cloudflare().SignUrl(ctx, file.FileKey, 3600*24) // 24 hours
			if err == nil {
				imageUrl = signedUrl
			} else {
				imageUrl = file.Url // Fallback
			}
		}
	}

	return &response.GenerationResponse{
		Status:    gen.Status,
		ImageUrl:  imageUrl,
		CreatedAt: gen.CreatedAt,
		UpdatedAt: gen.UpdatedAt,
	}, nil
}

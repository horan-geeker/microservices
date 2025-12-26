package aigc

import (
	"context"
	"microservices/logic/auth"

	entity "microservices/entity/model"
	"microservices/model"
	"microservices/service"

	"microservices/entity/consts"

	"google.golang.org/genai"
	"time"
)

type AIGCLogic interface {
	Generate(ctx context.Context, apiKey, model, prompt string) (uint64, error)
}

type aigcLogic struct {
	googleService service.Google
	model         model.Factory
	authLogic     auth.Logic
}

func NewAIGCLogic(googleService service.Google, model model.Factory) AIGCLogic {
	return &aigcLogic{
		googleService: googleService,
		model:         model,
	}
}

func (l *aigcLogic) Generate(ctx context.Context, apiKey, model, prompt string) (uint64, error) {
	auth, err := l.authLogic.GetAuthUser(ctx)
	if err != nil {
		return 0, err
	}
	// Create Generation Record
	gen := &entity.Generation{
		UserID:    auth.Uid,
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
	go func(genID uint64, prompt string) {
		resp, err := l.googleService.GenerateContent(
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

		_ = l.model.Generation().Update(ctx, genID, map[string]interface{}{
			"status":       status,
			"content_text": content,
			"updated_at":   time.Now(),
		})
	}(gen.ID, prompt)

	return gen.ID, nil
}

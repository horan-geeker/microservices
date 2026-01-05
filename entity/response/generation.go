package response

import (
	"microservices/entity/model"
	"time"
)

type GenerationResponse struct {
	Status    int       `json:"status"`
	ImageUrl  string    `json:"ImageUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GenerateResponse struct {
	GenerationID uint64 `json:"generationId"`
}

type GenerationSummary struct {
	ID            uint64      `json:"id"`
	Status        int         `json:"status"`
	OriginalFile  *model.File `json:"originalFile"`
	GeneratedFile *model.File `json:"generatedFile"`
	CreatedAt     time.Time   `json:"createdAt"`
}

type GetGenerationsResp struct {
	List  []*GenerationSummary `json:"list"`
	Total int64                `json:"total"`
}

type GenerationDetailResp struct {
	ID          uint64    `json:"id"`
	Status      int       `json:"status"`
	ContentText string    `json:"contentText"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

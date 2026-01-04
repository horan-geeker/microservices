package response

import "time"

type GenerationResponse struct {
	Status      int       `json:"status"`
	ContentText string    `json:"contentText"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GenerateResponse struct {
	GenerationID uint64 `json:"generationId"`
}

type GenerationSummary struct {
	ID        uint64    `json:"id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
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

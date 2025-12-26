package response

import "time"

type AIGCResponse struct {
	Status      int       `json:"status"`
	ContentText string    `json:"contentText"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type GenerateResponse struct {
	GenerationID uint64 `json:"generationId"`
}

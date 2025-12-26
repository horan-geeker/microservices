package request

type AIGCResultReq struct {
	ID uint64 `form:"id" binding:"required"`
}

type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

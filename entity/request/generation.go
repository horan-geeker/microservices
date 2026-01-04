package request

type GenerationResultReq struct {
	ID uint64 `form:"id" binding:"required"`
}

type GenerateRequest struct {
	Prompt string `json:"prompt"`
}

type GetGenerationsReq struct {
	Page int `form:"page" binding:"required,min=1"`
	Size int `form:"size" binding:"required,min=1,max=100"`
}

type GetGenerationDetailReq struct {
	ID uint64 `uri:"id" binding:"required"`
}

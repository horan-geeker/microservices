package controller

import (
	"microservices/cache"
	"microservices/entity/request"
	"microservices/entity/response"
	"microservices/logic"
	"microservices/logic/aigc"
	"microservices/model"
	"microservices/service"

	"github.com/gin-gonic/gin"
)

type Controller interface {
	Result(c *gin.Context, req *request.AIGCResultReq) (*response.AIGCResponse, error)
	Generate(c *gin.Context, req *request.GenerateRequest) (*response.GenerateResponse, error)
}

type controller struct {
	logic     logic.Factory
	aigcLogic aigc.AIGCLogic
	model     model.Factory
}

func NewAIGCController(model model.Factory, cache cache.Factory, service service.Factory) Controller {
	return &controller{
		logic:     logic.NewLogic(model, cache, service),
		aigcLogic: aigc.NewAIGCLogic(model, cache, service),
		model:     model,
	}
}

func (c *controller) Result(ctx *gin.Context, req *request.AIGCResultReq) (*response.AIGCResponse, error) {
	// Get Generation by ID
	generation, err := c.model.Generation().GetByID(ctx.Request.Context(), req.ID)
	if err != nil {
		return nil, err
	}

	return &response.AIGCResponse{
		Status:      generation.Status,
		ContentText: generation.ContentText,
		CreatedAt:   generation.CreatedAt,
		UpdatedAt:   generation.UpdatedAt,
	}, nil
}

func (c *controller) Generate(ctx *gin.Context, req *request.GenerateRequest) (*response.GenerateResponse, error) {
	// Call AIGC Logic
	id, err := c.aigcLogic.Generate(ctx.Request.Context(), "", "gemini-3-pro-preview", req.Prompt)
	if err != nil {
		return nil, err
	}
	return &response.GenerateResponse{
		GenerationID: id,
	}, nil
}

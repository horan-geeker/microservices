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
	Result(c *gin.Context, req *request.GenerationResultReq) (*response.GenerationResponse, error)
	Generate(c *gin.Context, req *request.GenerateRequest) (*response.GenerateResponse, error)
	List(c *gin.Context, req *request.GetGenerationsReq) (*response.GetGenerationsResp, error)
	Detail(c *gin.Context, id int) (*response.GenerationDetailResp, error)
}

type controller struct {
	logic     logic.Factory
	aigcLogic aigc.GenerationLogic
	model     model.Factory
}

func NewGenerationController(model model.Factory, cache cache.Factory, service service.Factory) Controller {
	return &controller{
		logic:     logic.NewLogic(model, cache, service),
		aigcLogic: aigc.NewGenerationLogic(model, cache, service),
		model:     model,
	}
}

func (c *controller) Result(ctx *gin.Context, req *request.GenerationResultReq) (*response.GenerationResponse, error) {
	auth, err := c.logic.Auth().GetAuthUser(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	return c.aigcLogic.GetResult(ctx.Request.Context(), auth.Uid, req.ID)
}

func (c *controller) Generate(ctx *gin.Context, req *request.GenerateRequest) (*response.GenerateResponse, error) {
	id, err := c.aigcLogic.Generate(ctx.Request.Context(), "", "gemini-3-pro-preview", req.Prompt)
	if err != nil {
		return nil, err
	}
	return &response.GenerateResponse{
		GenerationID: id,
	}, nil
}

func (c *controller) List(ctx *gin.Context, req *request.GetGenerationsReq) (*response.GetGenerationsResp, error) {
	auth, err := c.logic.Auth().GetAuthUser(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	return c.aigcLogic.List(ctx.Request.Context(), auth.Uid, req.Page, req.Size)
}

func (c *controller) Detail(ctx *gin.Context, id int) (*response.GenerationDetailResp, error) {
	auth, err := c.logic.Auth().GetAuthUser(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	return c.aigcLogic.Detail(ctx.Request.Context(), auth.Uid, uint64(id))
}

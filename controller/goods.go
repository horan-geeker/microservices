package controller

import (
	"microservices/cache"
	"microservices/entity/response"
	"microservices/logic"
	repo "microservices/model"
	"microservices/service"

	"github.com/gin-gonic/gin"
)

type GoodsController interface {
	GetList(c *gin.Context) (*response.GoodsListResponse, error)
}

type goodsController struct {
	logic logic.Factory
}

func NewGoodsController(model repo.Factory, cache cache.Factory, service service.Factory) GoodsController {
	return &goodsController{
		logic: logic.NewLogic(model, cache, service),
	}
}

// GetList godoc
// @Summary 获取商品列表
// @Description 获取所有上架的商品列表
// @Tags 商品
// @Accept json
// @Produce json
// @Success 200 {object} entity.Response[[]model.Goods]
// @Router /goods [get]
func (ctrl *goodsController) GetList(c *gin.Context) (*response.GoodsListResponse, error) {
	list, err := ctrl.logic.Goods().GetList(c.Request.Context())
	if err != nil {
		return nil, err
	}

	return &response.GoodsListResponse{
		List: list,
	}, nil
}

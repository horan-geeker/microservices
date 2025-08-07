package system

import (
	"github.com/gin-gonic/gin"
	"microservices/cache"
	"microservices/entity/response"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
)

type Controller interface {
	Health(c *gin.Context) (*response.Health, error)
}

type controller struct {
	logic logic.Factory
}

// Health .
func (s *controller) Health(c *gin.Context) (*response.Health, error) {
	return &response.Health{
		Status: "ok",
	}, nil
}

func NewController(model model.Factory, cache cache.Factory, service service.Factory) Controller {
	return &controller{
		logic: logic.NewLogic(model, cache, service),
	}
}

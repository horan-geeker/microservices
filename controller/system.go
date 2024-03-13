package controller

import (
	"github.com/gin-gonic/gin"
	"microservices/cache"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
)

type SystemApi interface {
	Health(c *gin.Context) (map[string]any, error)
}

type SystemController struct {
	logic logic.Factory
}

// Health .
func (s *SystemController) Health(c *gin.Context) (map[string]any, error) {
	return map[string]interface{}{
		"status": "UP",
	}, nil
}

func NewSystemController(model model.Factory, cache cache.Factory, service service.Factory) SystemApi {
	return &SystemController{
		logic: logic.NewLogic(model, cache, service),
	}
}

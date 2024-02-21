package controller

import (
	"github.com/gin-gonic/gin"
	"microservices/logic"
	"microservices/repository"
)

type SystemApi interface {
	Health(c *gin.Context) (map[string]any, error)
}

type SystemController struct {
	logic logic.LogicInterface
}

// Health .
func (s *SystemController) Health(c *gin.Context) (map[string]any, error) {
	return map[string]interface{}{
		"status": "UP",
	}, nil
}

func NewSystemController(repositoryFactory repository.Factory) SystemApi {
	return &SystemController{
		logic: logic.NewLogic(repositoryFactory),
	}
}

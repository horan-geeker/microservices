package system

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/entity"
)

// Health .
func Health(c *gin.Context) (*entity.Response, error) {
	return &entity.Response{
		Data: map[string]interface{}{
			"status": "UP",
		},
	}, nil
}

package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/entity"
)

// Login .
func Login(c *gin.Context) (entity.Response, error) {
	return entity.Response{
		Data: map[string]any{},
		Code: 0,
	}, nil
}

package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/entity"
	"microservices/internal/request"
)

// Login .
func Login(c *gin.Context) (*entity.Response, error) {
	param := request.LoginParams{}
	if err := c.ShouldBindJSON(&param); err != nil {
		return nil, err
	}
	return &entity.Response{
		Data: map[string]any{
			"username": param.Username,
		},
		Code: 0,
	}, nil
}

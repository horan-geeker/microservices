package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/entity"
)

func Logout(c *gin.Context) (entity.Response, error) {
	return entity.Response{}, nil
}

package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/request"
)

// Login .
func Login(c *gin.Context, req *request.LoginParams) (map[string]any, int, error) {
	return map[string]any{
		"username": req.Username,
	}, 0, nil
}

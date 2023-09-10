package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
)

// Login .
func Login(c *gin.Context, req *api.LoginParams) (map[string]any, int, error) {
	return map[string]any{
		"username": req.Username,
	}, 0, nil
}

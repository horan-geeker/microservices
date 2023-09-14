package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
)

// Login .
func (a *AuthController) Login(c *gin.Context, req *api.LoginParams) (map[string]any, int, error) {
	a.logic.Auth().AuthorizeCookie(c.Request.Context(), 1)
	return map[string]any{
		"username": req.Username,
	}, 0, nil
}

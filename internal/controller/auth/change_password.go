package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
)

type AuthController struct {
	authService *AuthService
}

// ChangePassword .
func ChangePassword(c *gin.Context, param *api.ChangePasswordParams) (map[string]any, int, error) {
	return map[string]any{
		"old_password": param.OldPassword,
	}, 0, nil
}

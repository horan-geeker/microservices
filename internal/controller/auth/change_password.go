package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
	"microservices/internal/pkg/ecode"
)

// ChangePassword .
func (a *AuthController) ChangePassword(c *gin.Context, param *api.ChangePasswordParams) (map[string]any, error) {
	if param.OldPassword == param.NewPassword {
		return nil, ecode.ErrUserPasswordDuplicate
	}
	return map[string]any{
		"old_password": param.OldPassword,
	}, nil
}

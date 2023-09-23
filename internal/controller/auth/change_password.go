package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
)

// ChangePassword .
func (a *AuthController) ChangePassword(c *gin.Context, param *api.ChangePasswordParams) (map[string]any, error) {
	if err := a.logic.Auth().ChangePassword(c.Request.Context(), param.Uid, param.NewPassword, param.OldPassword, param.SmsCode); err != nil {
		return nil, err
	}
	return nil, nil
}

package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
)

func (a *AuthController) ChangePasswordByPhone(c *gin.Context, param api.ChangePasswordByPhoneParams) (map[string]any, error) {
	if err := a.logic.Auth().ChangePasswordByPhone(c.Request.Context(), param.NewPassword, param.Phone, param.SmsCode); err != nil {
		return nil, err
	}
	return nil, nil
}

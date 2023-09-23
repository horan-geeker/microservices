package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
	"time"
)

// Login .
func (a *AuthController) Login(c *gin.Context, params *api.LoginParams) (map[string]any, error) {
	user, err := a.logic.Auth().Login(c.Request.Context(), params.Name, params.Email, params.Phone,
		params.Password, params.SmsCode)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"user": map[string]any{
			"id":          user.ID,
			"status":      user.Status,
			"name":        user.Name,
			"email":       user.Email,
			"phone":       user.Phone,
			"lastLoginAt": user.LoginAt.Format(time.DateTime),
		},
	}, nil
}

package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
)

// ChangePassword .
func (a *AuthController) ChangePassword(c *gin.Context, param *api.ChangePasswordParams) (map[string]any, error) {
	auth, err := a.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	if err := a.logic.Auth().ChangePassword(c.Request.Context(), auth.Uid, param.NewPassword, param.OldPassword); err != nil {
		return nil, err
	}
	return nil, nil
}

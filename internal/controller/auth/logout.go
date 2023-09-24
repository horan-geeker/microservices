package auth

import (
	"github.com/gin-gonic/gin"
)

// Logout .
func (a *AuthController) Logout(c *gin.Context) (map[string]any, error) {
	auth, err := a.logic.Auth().GetAuthUser(c.Request.Context())
	if err != nil {
		return nil, err
	}
	if err := a.logic.Auth().Logout(c.Request.Context(), auth.Uid); err != nil {
		return nil, err
	}
	return nil, nil
}

package auth

import (
	"github.com/gin-gonic/gin"
)

// Logout .
func (a *AuthController) Logout(c *gin.Context) (map[string]any, int, error) {
	return map[string]any{}, 1, nil
}

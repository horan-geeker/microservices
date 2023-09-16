package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
)

// Logout .
func (a *AuthController) Logout(c *gin.Context) (map[string]any, error) {
	return map[string]any{}, errors.New("foo")
}

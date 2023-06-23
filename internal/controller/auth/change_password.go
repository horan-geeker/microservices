package auth

import (
	"github.com/gin-gonic/gin"
	"microservices/internal/request"
)

// ChangePassword .
func ChangePassword(c *gin.Context, param *request.ChangePasswordParams) (map[string]any, int, error) {
	return map[string]any{
		"old_password": param.OldPassword,
	}, 0, nil
}

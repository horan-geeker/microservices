package user

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
)

// Register .
func (a *UserController) Register(c *gin.Context, params *api.RegisterParams) (map[string]any, error) {
	user, token, err := a.logic.Auth().Register(c.Request.Context(), params.Name, params.Email, params.Phone, params.Password)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"user": map[string]any{
			"id": user.ID,
		},
		"token": token,
	}, nil
}

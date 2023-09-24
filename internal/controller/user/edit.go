package user

import (
	"github.com/gin-gonic/gin"
	"microservices/api"
	"microservices/internal/pkg/ecode"
)

func (u *UserController) Edit(c *gin.Context, param *api.EditUserParam) (map[string]any, error) {
	token := c.GetHeader("Authorization")
	if len(token) == 0 {
		return nil, ecode.ErrTokenIsEmpty
	}
	auth, err := u.logic.Auth().GetAuthUser(token)
	if err != nil {
		return nil, err
	}
	err = u.logic.Users().Edit(c.Request.Context(), auth.Uid, param.Name, param.Email, param.Phone)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

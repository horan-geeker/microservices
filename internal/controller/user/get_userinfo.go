package user

import (
	"github.com/gin-gonic/gin"
)

// GetUserinfo .
func (u *UserController) GetUserinfo(c *gin.Context, uid int) (map[string]any, error) {
	userinfo, err := u.logic.Users().GetByUid(c.Request.Context(), uid)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"user": userinfo,
	}, nil
}

package user

import (
	"github.com/gin-gonic/gin"
)

// GetUserinfo .
func (u *UserController) GetUserinfo(c *gin.Context, uid int) (map[string]any, int, error) {
	c.Param("id")
	userinfo, err := u.logic.Users().GetByUid(c.Request.Context(), uid)
	if err != nil {
		return nil, 0, err
	}
	return map[string]any{
		"user": userinfo,
	}, 0, nil
}

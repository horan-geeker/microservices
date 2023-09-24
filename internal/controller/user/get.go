package user

import (
	"github.com/gin-gonic/gin"
)

// Get .
func (u *UserController) Get(c *gin.Context, uid uint64) (map[string]any, error) {
	userinfo, err := u.logic.Users().GetByUid(c.Request.Context(), uid)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"user": userinfo,
	}, nil
}

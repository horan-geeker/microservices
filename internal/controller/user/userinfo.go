package user

import "github.com/gin-gonic/gin"

// Userinfo .
func (u *UserController) Userinfo(c *gin.Context) (map[string]any, int, error) {
	userinfo, err := u.logic.Users().GetByUserName(c.Request.Context(), "")
	if err != nil {
		return nil, 0, err
	}
	return map[string]any{
		"user": userinfo,
	}, 0, nil
}

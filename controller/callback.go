package controller

import (
	"microservices/entity/ecode"
	"microservices/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CallbackController struct {
	logic logic.Factory
}

func NewCallbackController(l logic.Factory) *CallbackController {
	return &CallbackController{logic: l}
}

// GoogleAuthCallback handles the Google OAuth callback.
func (c *CallbackController) GoogleAuthCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": ecode.ErrParamInvalid.Code, "msg": "code is required", "data": nil})
		return
	}

	user, token, err := c.logic.Callback().GoogleCallback(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": err.Error(), "data": nil})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})
}

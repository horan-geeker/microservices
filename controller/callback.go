package controller

import (
	"microservices/entity/request"
	"microservices/entity/response"
	"microservices/logic"

	"github.com/gin-gonic/gin"
)

type CallbackController interface {
	GoogleAuthCallback(ctx *gin.Context, param *request.GoogleCallbackParam) (*response.GoogleCallback, error)
}

type callbackController struct {
	logic logic.Factory
}

func NewCallbackController(l logic.Factory) CallbackController {
	return &callbackController{logic: l}
}

// GoogleAuthCallback handles the Google OAuth callback.
// @Summary Google OAuth Callback
// @Description Handle Google OAuth callback code
// @Tags callback
// @Accept  json
// @Produce  json
// @Param code query string true "OAuth Code"
// @Success 200 {object} entity.Response[response.GoogleCallback]
// @Failure 400 {object} entity.Response[any]
// @Router /callback/google-auth [get]
func (c *callbackController) GoogleAuthCallback(ctx *gin.Context, param *request.GoogleCallbackParam) (*response.GoogleCallback, error) {
	user, token, err := c.logic.Callback().GoogleCallback(ctx, param.Code)
	if err != nil {
		return nil, err
	}

	return &response.GoogleCallback{
		User:  user,
		Token: token,
	}, nil
}

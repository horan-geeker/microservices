package controller

import (
	"github.com/gin-gonic/gin"
	"microservices/cache"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
)

type NotifyApi interface {
	SendSms(c *gin.Context) (map[string]any, error)
	SendEmail(c *gin.Context) (map[string]any, error)
}

type notifyController struct {
	logic logic.Factory
}

func (n *notifyController) SendEmail(c *gin.Context) (map[string]any, error) {
	//TODO implement me
	panic("implement me")
}

func (n *notifyController) SendSms(c *gin.Context) (map[string]any, error) {
	if err := n.logic.Notify().SendSmsCode(c.Request.Context(), "13571899655", "1234"); err != nil {
		return nil, err
	}
	return nil, nil
}

func NewNotifyController(model model.Factory, cache cache.Factory, service service.Factory) NotifyApi {
	return &notifyController{
		logic: logic.NewLogic(model, cache, service),
	}
}

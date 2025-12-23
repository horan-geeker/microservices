package controller

import (
	"github.com/gin-gonic/gin"
	"microservices/cache"
	"microservices/entity/response"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
)

type NotifyController interface {
	SendSms(c *gin.Context) (*response.SendSms, error)
	SendEmail(c *gin.Context) (*response.SendEmail, error)
}

type notifyController struct {
	logic logic.Factory
}

func (n *notifyController) SendSms(c *gin.Context) (*response.SendSms, error) {
	if err := n.logic.Notify().SendSmsCode(c.Request.Context(), "13571899655", "1234"); err != nil {
		return nil, err
	}
	return nil, nil
}

func (n *notifyController) SendEmail(c *gin.Context) (*response.SendEmail, error) {
	//TODO implement me
	panic("implement me")
}

func NewNotifyController(model model.Factory, cache cache.Factory, service service.Factory) NotifyController {
	return &notifyController{
		logic: logic.NewLogic(model, cache, service),
	}
}

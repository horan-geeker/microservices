package controller

import (
	"github.com/gin-gonic/gin"
	"microservices/logic"
	"microservices/repository"
)

type NotifyApi interface {
	SendSms(c *gin.Context) (map[string]any, error)
	SendEmail(c *gin.Context) (map[string]any, error)
}

type notifyController struct {
	logic logic.LogicInterface
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

func NewNotifyController(factory repository.Factory) NotifyApi {
	return &notifyController{
		logic: logic.NewLogic(factory),
	}
}

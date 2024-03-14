package notify

import (
	"github.com/gin-gonic/gin"
	"microservices/cache"
	"microservices/logic"
	"microservices/model"
	"microservices/service"
)

type Controller interface {
	SendSms(c *gin.Context) (map[string]any, error)
	SendEmail(c *gin.Context) (map[string]any, error)
}

type controller struct {
	logic logic.Factory
}

func (n *controller) SendEmail(c *gin.Context) (map[string]any, error) {
	//TODO implement me
	panic("implement me")
}

func (n *controller) SendSms(c *gin.Context) (map[string]any, error) {
	if err := n.logic.Notify().SendSmsCode(c.Request.Context(), "13571899655", "1234"); err != nil {
		return nil, err
	}
	return nil, nil
}

func NewController(model model.Factory, cache cache.Factory, service service.Factory) Controller {
	return &controller{
		logic: logic.NewLogic(model, cache, service),
	}
}

package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/service"
)

func init() {
	notifyController := controller.NewNotifyController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.POST("/notify/sms", notifyController.SendSms)
}

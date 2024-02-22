package router

import (
	"microservices/controller"
	"microservices/repository"
	"microservices/service"
)

func init() {
	notifyController := controller.NewNotifyController(repository.NewFactory(), service.NewFactory())
	router.POST("/notify/sms", notifyController.SendSms)
}

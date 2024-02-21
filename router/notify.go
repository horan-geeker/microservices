package router

import (
	"microservices/controller"
	"microservices/repository"
)

func init() {
	notifyController := controller.NewNotifyController(repository.NewFactory())
	router.POST("/notify/sms", notifyController.SendSms)
}

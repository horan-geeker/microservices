package router

import (
	"microservices/controller"
	"microservices/pkg/app"
	"microservices/repository"
)

func init() {
	notifyController := controller.NewNotifyController(repository.NewFactory())
	router := app.GetApp()
	router.POST("/notify/sms", notifyController.SendSms)
}

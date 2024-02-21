package router

import (
	"microservices/controller"
)

func init() {
	controller := controller.NewNotifyController(repositoryFactory)
	router.POST("/notify/sms", controller.SendSms)
}

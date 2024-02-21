package router

import (
	"microservices/service"
)

func init() {
	controller := service.NewNotifyController(repositoryFactory)
	router.POST("/notify/sms", controller.SendSms)
}

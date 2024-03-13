package router

import (
	"microservices/cache"
	"microservices/controller/notify"
	"microservices/model"
	"microservices/service"
)

func init() {
	notifyController := notify.NewController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.POST("/notify/sms", notifyController.SendSms)
}

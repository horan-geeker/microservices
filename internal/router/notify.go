package router

import (
	"microservices/internal/controller/notify"
)

func init() {
	controller := notify.NewNotifyController(dataFactory, cacheFactory, serviceFactory)
	router.POST("/notify/sms", controller.SendSms)
}

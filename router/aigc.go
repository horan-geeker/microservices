package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/service"
)

func init() {
	aigcController := controller.NewAIGCController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.POST("/aigc/generate", aigcController.Generate)
	router.GET("/aigc/result", aigcController.Result)
}

package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	aigcController := controller.NewAIGCController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.POST("/aigc/generate", middleware.Authenticate(), aigcController.Generate)
	router.GET("/aigc/result", middleware.Authenticate(), aigcController.Result)
}

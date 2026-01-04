package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/router/middleware"
	"microservices/service"
)

func init() {
	generationController := controller.NewGenerationController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.POST("/generate", middleware.Authenticate(), generationController.Generate)
	router.GET("/generate/result", middleware.Authenticate(), generationController.Result)
	router.GET("/generations", middleware.Authenticate(), generationController.List)
	router.GET("/generations/:id", middleware.Authenticate(), generationController.Detail)
}

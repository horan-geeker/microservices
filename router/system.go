package router

import (
	"microservices/cache"
	"microservices/controller"
	"microservices/model"
	"microservices/service"
)

func init() {
	systemController := controller.NewSystemController(model.NewFactory(), cache.NewFactory(), service.NewFactory())
	router.GET("/system/health", systemController.Health)
}

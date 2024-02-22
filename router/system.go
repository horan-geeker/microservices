package router

import (
	"microservices/controller"
	"microservices/repository"
	"microservices/service"
)

func init() {
	systemController := controller.NewSystemController(repository.NewFactory(), service.NewFactory())
	router.GET("/system/health", systemController.Health)
}

package router

import (
	"microservices/controller"
	"microservices/repository"
)

func init() {
	systemController := controller.NewSystemController(repository.NewFactory())
	router.GET("/system/health", systemController.Health)
}

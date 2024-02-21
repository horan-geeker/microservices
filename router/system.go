package router

import (
	"microservices/controller"
	"microservices/pkg/app"
	"microservices/repository"
)

func init() {
	systemController := controller.NewSystemController(repository.NewFactory())
	router := app.GetApp()
	router.GET("/system/health", systemController.Health)
}

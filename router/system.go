package router

import (
	"microservices/controller"
)

func init() {
	systemController := controller.NewSystemController(repositoryFactory)
	router.GET("/system/health", systemController.Health)
}

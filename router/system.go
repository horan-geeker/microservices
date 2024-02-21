package router

import (
	"microservices/service"
)

func init() {
	systemController := service.NewSystemController(repositoryFactory)
	router.GET("/system/health", systemController.Health)
}

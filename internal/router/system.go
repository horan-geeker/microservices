package router

import (
	"microservices/internal/controller/system"
)

func init() {
	systemController := system.NewSystemController(dataFactory, cacheFactory, serviceFactory)
	router.GET("/system/health", systemController.Health)
}

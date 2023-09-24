package router

import (
	"microservices/internal/controller/system"
)

func init() {
	router.GET("/system/health", system.Health)
}

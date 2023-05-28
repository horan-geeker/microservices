package router

import (
	"microservices/internal/controller/system"
	"net/http"
)

func init() {
	routes = append(routes,
		router{
			Method: http.MethodGet,
			Path:   "/system/health",
			Func:   system.Health,
		},
	)
}

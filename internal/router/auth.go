package router

import (
	"microservices/internal/controller/auth"
	"microservices/internal/request"
	"net/http"
)

func init() {
	routes = append(routes,
		router{
			Method:  http.MethodPost,
			Path:    "/auth/login",
			Request: request.LoginParams{},
			Func:    auth.Login,
		},
		router{
			Method: http.MethodPost,
			Path:   "/auth/logout",
			Func:   auth.Logout,
		},
	)
}

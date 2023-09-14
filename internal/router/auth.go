package router

import (
	"fmt"
	"microservices/internal/controller/auth"
	"net/http"
)

func init() {
	fmt.Println("auth init")
	authController := auth.NewAuthController(dataFactory, cacheFactory)
	routes = append(routes, []router{
		{
			Method: http.MethodPost,
			Path:   "/auth/login",
			Func:   authController.Login,
		},
		{
			Method: http.MethodPost,
			Path:   "/auth/logout",
			Func:   authController.Logout,
		},
		{
			Method: http.MethodPost,
			Path:   "/auth/change-password",
			Func:   authController.ChangePassword,
		},
	}...,
	)
}

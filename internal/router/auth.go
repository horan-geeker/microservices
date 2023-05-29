package router

import (
	"microservices/internal/controller/auth"
	"net/http"
)

func init() {
	routes = append(routes, []router{
		{
			Method: http.MethodPost,
			Path:   "/auth/login",
			Func:   auth.Login,
		},
		{
			Method: http.MethodPost,
			Path:   "/auth/logout",
			Func:   auth.Logout,
		},
		{
			Method: http.MethodPost,
			Path:   "/auth/change-password",
			Func:   auth.ChangePassword,
		},
	}...,
	)
}

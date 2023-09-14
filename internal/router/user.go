package router

import (
	"microservices/internal/controller/user"
	"net/http"
)

func init() {
	userController := user.NewUserController(dataFactory, cacheFactory)
	routes = append(routes, []router{
		{
			Method: http.MethodGet,
			Path:   "/user/:uid",
			Func:   userController.GetUserinfo,
		},
	}...,
	)
}

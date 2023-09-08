package router

import (
	"microservices/internal/controller/user"
	"microservices/internal/store/mysql"
	"net/http"
)

func init() {
	store := mysql.ConnectDB(nil)
	userController := user.NewUserController(store)
	routes = append(routes, []router{
		{
			Method: http.MethodGet,
			Path:   "/userinfo/",
			Func:   userController.Userinfo,
		},
	}...,
	)
}

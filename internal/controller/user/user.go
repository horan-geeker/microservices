package user

import (
	"microservices/internal/logic"
	"microservices/internal/store"
)

type UserController struct {
	logic logic.LogicInterface
}

// NewUserController .
func NewUserController(store store.DataFactory, cache store.CacheFactory) *UserController {
	return &UserController{
		logic: logic.NewLogic(store, cache),
	}
}

package logic

import (
	"microservices/internal/service"
	"microservices/internal/store"
)

// LogicInterface defines functions used to return resource interface.
type LogicInterface interface {
	Users() UserLogicInterface
	Auth() AuthLogicInterface
	Notify() NotifyLogicInterface
}

type logic struct {
	store store.Factory
	cache store.CacheFactory
	srv   service.Factory
}

func (l *logic) Notify() NotifyLogicInterface {
	return newNotify(l)
}

func (l *logic) Users() UserLogicInterface {
	return newUsers(l)
}

func (l *logic) Auth() AuthLogicInterface {
	return newAuth(l)
}

// NewLogic .
func NewLogic(store store.Factory, cache store.CacheFactory, srv service.Factory) LogicInterface {
	return &logic{store: store, cache: cache, srv: srv}
}

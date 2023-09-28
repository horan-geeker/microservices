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
	store store.DataFactory
	cache store.CacheFactory
	srv   service.ServiceFactory
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
func NewLogic(store store.DataFactory, cache store.CacheFactory, srv service.ServiceFactory) LogicInterface {
	return &logic{store: store, cache: cache, srv: srv}
}

package logic

import "microservices/internal/store"

// LogicInterface defines functions used to return resource interface.
type LogicInterface interface {
	Users() UserLogicInterface
	Auth() AuthLogicInterface
}

type logic struct {
	store store.DataFactory
	cache store.CacheFactory
}

func (l *logic) Users() UserLogicInterface {
	return newUsers(l)
}

func (l *logic) Auth() AuthLogicInterface {
	return newAuth(l)
}

// NewLogic .
func NewLogic(store store.DataFactory, cache store.CacheFactory) LogicInterface {
	return &logic{store: store, cache: cache}
}

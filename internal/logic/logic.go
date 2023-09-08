package logic

import "microservices/internal/store"

// LogicInterface defines functions used to return resource interface.
type LogicInterface interface {
	Users() UserLogicInterface
	Auth() AuthLogicInterface
}

type logic struct {
	store store.Factory
}

func (l *logic) Users() UserLogicInterface {
	return newUsers(l)
}

func (l *logic) Auth() AuthLogicInterface {
	return newAuth(l)
}

// NewLogic .
func NewLogic(store store.Factory) LogicInterface {
	return &logic{store: store}
}

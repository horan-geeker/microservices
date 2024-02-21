package logic

import (
	"microservices/logic/auth"
	"microservices/logic/notify"
	"microservices/logic/user"
	"microservices/repository"
)

// LogicInterface defines functions used to return resource interface.
type LogicInterface interface {
	Users() user.UserLogicInterface
	Auth() auth.AuthLogicInterface
	Notify() notify.NotifyLogicInterface
}

type logic struct {
	repository repository.Factory
}

func (l *logic) Notify() notify.NotifyLogicInterface {
	return notify.NewNotify(l.repository)
}

func (l *logic) Users() user.UserLogicInterface {
	return user.NewUsers(l.repository)
}

func (l *logic) Auth() auth.AuthLogicInterface {
	return auth.NewAuth(l.repository)
}

// NewLogic .
func NewLogic(factory repository.Factory) LogicInterface {
	return &logic{factory}
}

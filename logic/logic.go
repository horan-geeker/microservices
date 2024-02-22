package logic

import (
	"microservices/logic/auth"
	"microservices/logic/notify"
	"microservices/logic/user"
	"microservices/repository"
	"microservices/service"
)

// LogicInterface defines functions used to return resource interface.
type LogicInterface interface {
	Users() user.UserLogicInterface
	Auth() auth.AuthLogicInterface
	Notify() notify.NotifyLogicInterface
}

type logic struct {
	repository repository.Factory
	service    service.Factory
}

func (l *logic) Notify() notify.NotifyLogicInterface {
	return notify.NewNotify(l.repository, l.service)
}

func (l *logic) Users() user.UserLogicInterface {
	return user.NewUsers(l.repository, l.service)
}

func (l *logic) Auth() auth.AuthLogicInterface {
	return auth.NewAuth(l.repository, l.service)
}

// NewLogic .
func NewLogic(repo repository.Factory, service service.Factory) LogicInterface {
	return &logic{repo, service}
}

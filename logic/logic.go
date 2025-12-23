package logic

import (
	"microservices/cache"
	"microservices/logic/auth"
	"microservices/logic/callback"
	"microservices/logic/notify"
	"microservices/logic/user"
	"microservices/model"
	"microservices/service"
)

// Factory defines functions used to return resource interface.
type Factory interface {
	User() user.Logic
	Auth() auth.Logic
	Callback() callback.Logic
	Notify() notify.Logic
}

type factory struct {
	model   model.Factory
	cache   cache.Factory
	service service.Factory
}

func (l *factory) Notify() notify.Logic {
	return notify.NewNotify(l.model, l.cache, l.service)
}

func (l *factory) User() user.Logic {
	return user.NewLogic(l.model, l.cache, l.service)
}

func (l *factory) Auth() auth.Logic {
	return auth.NewAuth(l.model, l.cache, l.service)
}

func (l *factory) Callback() callback.Logic {
	return callback.NewCallback(l.model, l.cache, l.service)
}

// NewLogic .
func NewLogic(model model.Factory, cache cache.Factory, service service.Factory) Factory {
	return &factory{model, cache, service}
}

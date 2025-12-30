package goods

import (
	"context"
	"microservices/cache"
	"microservices/entity/model"
	repo "microservices/model"
	"microservices/service"
)

type Logic interface {
	GetList(ctx context.Context) ([]*model.Goods, error)
}

type logic struct {
	model repo.Factory
	cache cache.Factory
	srv   service.Factory
}

func NewLogic(model repo.Factory, cache cache.Factory, service service.Factory) Logic {
	return &logic{model: model, cache: cache, srv: service}
}

func (l *logic) GetList(ctx context.Context) ([]*model.Goods, error) {
	return l.model.Goods().GetList(ctx)
}

package notify

import (
	"context"
	"microservices/cache"
	"microservices/model"
	"microservices/service"
)

type Logic interface {
	SendSmsCode(ctx context.Context, phone string, code string) error
}

type logic struct {
	model model.Factory
	cache cache.Factory
	srv   service.Factory
}

func (n *logic) SendSmsCode(ctx context.Context, phone string, code string) error {
	if err := n.cache.Auth().SetSmsCode(ctx, phone, code); err != nil {

	}
	return n.srv.Aliyun().SendSMSCode(ctx, phone, code)
}

func NewNotify(model model.Factory, cache cache.Factory, service service.Factory) Logic {
	return &logic{
		model: model,
		cache: cache,
		srv:   service,
	}
}

package logic

import (
	"context"
	"microservices/internal/service"
	"microservices/internal/store"
)

type NotifyLogicInterface interface {
	SendSmsCode(ctx context.Context, phone string, code string) error
}

type notifyLogic struct {
	cache store.CacheFactory
	srv   service.ServiceFactory
}

func (n *notifyLogic) SendSmsCode(ctx context.Context, phone string, code string) error {
	if err := n.cache.Auth().SetSmsCode(ctx, phone, code); err != nil {

	}
	return n.srv.Aliyun().SendSMSCode(ctx, phone, code)
}

func newNotify(logic *logic) NotifyLogicInterface {
	return &notifyLogic{
		cache: logic.cache,
		srv:   logic.srv,
	}
}

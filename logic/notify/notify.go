package notify

import (
	"context"
	"microservices/repository"
	"microservices/service"
)

type NotifyLogicInterface interface {
	SendSmsCode(ctx context.Context, phone string, code string) error
}

type notifyLogic struct {
	repo repository.Factory
	srv  service.Factory
}

func (n *notifyLogic) SendSmsCode(ctx context.Context, phone string, code string) error {
	if err := n.repo.Auth().SetSmsCode(ctx, phone, code); err != nil {

	}
	return n.srv.Aliyun().SendSMSCode(ctx, phone, code)
}

func NewNotify(repo repository.Factory, service service.Factory) NotifyLogicInterface {
	return &notifyLogic{
		repo: repo,
		srv:  service,
	}
}

package notify

import (
	"context"
	"microservices/repository"
)

type NotifyLogicInterface interface {
	SendSmsCode(ctx context.Context, phone string, code string) error
}

type notifyLogic struct {
	repo repository.Factory
}

func (n *notifyLogic) SendSmsCode(ctx context.Context, phone string, code string) error {
	if err := n.repo.Auth().SetSmsCode(ctx, phone, code); err != nil {

	}
	return n.repo.Aliyun().SendSMSCode(ctx, phone, code)
}

func NewNotify(repo repository.Factory) NotifyLogicInterface {
	return &notifyLogic{
		repo: repo,
	}
}

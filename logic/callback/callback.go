package callback

import (
	"context"
	"github.com/go-pay/gopay"
	"microservices/cache"
	entity "microservices/entity/model"
	"microservices/entity/request"
	"microservices/model"
	"microservices/service"
)

type Logic interface {
	GoogleCallback(ctx context.Context, code string) (*entity.User, string, error)
	HandleAlipayNotify(ctx context.Context, payload map[string]string) error
	HandleAlipayCallback(ctx context.Context, params map[string]string) error
	HandleAppleCallback(ctx context.Context, notification *request.AppleIAPNotification) error
	VerifyAlipayNotifySign(ctx context.Context, bm gopay.BodyMap) error
	HandleStripeCallback(ctx context.Context, payload []byte, header string) error
}

type logic struct {
	model model.Factory
	cache cache.Factory
	srv   service.Factory
}

func NewCallback(model model.Factory, cache cache.Factory, service service.Factory) Logic {
	return &logic{
		model: model,
		cache: cache,
		srv:   service,
	}
}

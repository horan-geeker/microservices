package order

import (
	"context"
	"fmt"
	"microservices/entity/consts"
	entity "microservices/entity/model"
	"microservices/pkg/util"
	"time"
)

func (p *logic) CreateStripeCheckout(ctx context.Context, userId int, priceId, clientIP string) (*entity.Order, string, error) {
	var outTradeNo string
	outTradeNo = fmt.Sprintf("%s%d%010d", time.Now().Format(consts.TimeFormatSeq), util.RandomN(3), userId)
	expiredAt := time.Now().Add(30 * time.Minute)

	goods, err := p.model.Goods().GetByPriceId(ctx, priceId)
	if err != nil {
		return nil, "", err
	}
	// Create order record
	order := &entity.Order{
		OutTradeNo:  outTradeNo,
		UserId:      userId,
		Status:      consts.OrderStatusPending,
		TradeType:   "stripe",
		Platform:    "",
		ClientIp:    &clientIP,
		ExpireAt:    &expiredAt,
		Subject:     goods.Name,
		Description: &goods.Description,
		TotalAmount: uint64(goods.Price),
		Currency:    goods.Currency,
	}

	if err := p.model.Order().Create(ctx, order); err != nil {
		return nil, "", fmt.Errorf("Creating order failed: %w", err)
	}
	user, err := p.model.User().GetByUid(ctx, userId)
	if err != nil {
		return nil, "", err
	}
	url, err := p.srv.Stripe().CreateCheckoutSession(ctx, user, priceId, order.ID)
	if err != nil {
		_ = p.model.Order().Update(ctx, order.ID, map[string]interface{}{
			"status": consts.OrderStatusCreateError,
		})
		return nil, "", err
	}

	return order, url, nil
}

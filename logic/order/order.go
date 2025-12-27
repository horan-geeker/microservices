package order

import (
	"context"
	"microservices/cache"
	entity "microservices/entity/model"
	"microservices/entity/request"
	"microservices/entity/response"
	"microservices/model"
	"microservices/service"
)

// Logic defines functions used to handle payment api.
type Logic interface {
	CreateAliPayPrepayOrder(ctx context.Context, userId, price int, platform, channel, clientIP,
		description string) (*entity.Order, string, error)
	VerifyAppleReceipt(ctx context.Context, userId int, receiptData string,
		excludeOldTransactions bool) (*response.AppleVerifyReceipt, error)
	GetDetail(ctx context.Context, id, userId int) (*entity.Order, error)
	GetList(ctx context.Context, userID int, req *request.GetOrderListRequest) (*response.GetOrderListResponse, error)
}

type logic struct {
	model model.Factory
	cache cache.Factory
	srv   service.Factory
}

func NewLogic(model model.Factory, cache cache.Factory, service service.Factory) Logic {
	return &logic{model: model, cache: cache, srv: service}
}

func (l *logic) GetDetail(ctx context.Context, id, userId int) (*entity.Order, error) {
	return l.model.Order().GetByIdAndUserId(ctx, id, userId)
}

func (l *logic) GetList(ctx context.Context, userID int, req *request.GetOrderListRequest) (*response.GetOrderListResponse, error) {
	total, err := l.model.Order().CountByUserId(ctx, uint64(userID))
	if err != nil {
		return nil, err
	}

	offset := (req.Page - 1) * req.PageSize
	orders, err := l.model.Order().GetByUserId(ctx, uint64(userID), req.PageSize, offset)
	if err != nil {
		return nil, err
	}

	return &response.GetOrderListResponse{
		Items: orders,
		Total: total,
	}, nil
}

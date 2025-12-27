package model

import (
	"context"
	"microservices/entity/model"
	"time"

	"gorm.io/gorm"
)

// Order defines the order storage interface.
type Order interface {
	Create(ctx context.Context, order *model.Order) error
	Update(ctx context.Context, id int, data map[string]interface{}) error
	UpdateByOutTradeNo(ctx context.Context, mchId, outTradeNo string, data map[string]interface{}) error
	GetByIdAndUserId(ctx context.Context, id, userId int) (*model.Order, error)
	GetByOutTradeNo(ctx context.Context, outTradeNo string) (*model.Order, error)
	GetByUserId(ctx context.Context, userId uint64, limit, offset int) ([]*model.Order, error)
	CountByUserId(ctx context.Context, userId uint64) (int64, error)
}

type order struct {
	db *gorm.DB
}

func newOrder(s *factory) Order {
	return &order{
		db: s.db.Model(&model.Order{}),
	}
}

// Create .
func (o *order) Create(ctx context.Context, order *model.Order) error {
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	return o.db.WithContext(ctx).Create(order).Error
}

// Update .
func (o *order) Update(ctx context.Context, id int, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return o.db.WithContext(ctx).Where("id = ?", id).Updates(data).Error
}

// UpdateByOutTradeNo .
func (o *order) UpdateByOutTradeNo(ctx context.Context, mchId, outTradeNo string, data map[string]interface{}) error {
	data["updated_at"] = time.Now()
	return o.db.WithContext(ctx).Where("mchid = ? AND out_trade_no = ?", mchId, outTradeNo).Updates(data).Error
}

// GetByIdAndUserId .
func (o *order) GetByIdAndUserId(ctx context.Context, id, userId int) (t *model.Order, err error) {
	err = o.db.WithContext(ctx).Where("id = ?", id).Where("user_id", userId).Take(&t).Error
	return
}

// GetByOutTradeNo .
func (o *order) GetByOutTradeNo(ctx context.Context, outTradeNo string) (t *model.Order, err error) {
	err = o.db.WithContext(ctx).Where("out_trade_no = ?", outTradeNo).Take(&t).Error
	return
}

// GetByUserId .
func (o *order) GetByUserId(ctx context.Context, userId uint64, limit, offset int) (orders []*model.Order, err error) {
	err = o.db.WithContext(ctx).Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return
}

// CountByUserId .
func (o *order) CountByUserId(ctx context.Context, userId uint64) (count int64, err error) {
	err = o.db.WithContext(ctx).Where("user_id = ?", userId).Count(&count).Error
	return
}

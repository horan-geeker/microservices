package model

import (
	"context"
	"microservices/entity/model"

	"gorm.io/gorm"
)

// Goods defines the goods storage interface.
type Goods interface {
	GetList(ctx context.Context) ([]*model.Goods, error)
	GetByPriceId(ctx context.Context, priceId string) (model.Goods, error)
}

type goods struct {
	db *gorm.DB
}

func newGoods(s *factory) Goods {
	return &goods{
		db: s.db.Model(&model.Goods{}),
	}
}

// GetList returns all goods.
func (g *goods) GetList(ctx context.Context) (list []*model.Goods, err error) {
	err = g.db.WithContext(ctx).Where("status = ?", 1).Find(&list).Error
	return
}

func (g *goods) GetByPriceId(ctx context.Context, priceId string) (t model.Goods, err error) {
	err = g.db.WithContext(ctx).Where("stripe_price_id = ?", priceId).Take(&t).Error
	return
}

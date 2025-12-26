package model

import (
	"context"
	entity "microservices/entity/model"

	"gorm.io/gorm"
)

type Generation interface {
	Create(ctx context.Context, generation *entity.Generation) error
	Update(ctx context.Context, id uint64, updates map[string]interface{}) error
	GetByID(ctx context.Context, id uint64) (*entity.Generation, error)
}

type generation struct {
	db *gorm.DB
}

func newGeneration(f *factory) Generation {
	return &generation{
		db: f.db,
	}
}

func (g *generation) Create(ctx context.Context, generation *entity.Generation) error {
	return g.db.WithContext(ctx).Create(generation).Error
}

func (g *generation) Update(ctx context.Context, id uint64, updates map[string]interface{}) error {
	return g.db.WithContext(ctx).Model(&entity.Generation{}).Where("id = ?", id).Updates(updates).Error
}

func (g *generation) GetByID(ctx context.Context, id uint64) (*entity.Generation, error) {
	var gen entity.Generation
	err := g.db.WithContext(ctx).First(&gen, id).Error
	if err != nil {
		return nil, err
	}
	return &gen, nil
}

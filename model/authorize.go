package model

import (
	"context"
	"microservices/entity/model"
	"time"

	"gorm.io/gorm"
)

// Authorize defines the authorize storage interface.
type Authorize interface {
	Create(ctx context.Context, authorize *model.Authorize) error
	Update(ctx context.Context, id int, data map[string]any) error
	GetByProvider(ctx context.Context, provider, providerID string) (*model.Authorize, error)
	GetByUserID(ctx context.Context, userID int) ([]*model.Authorize, error)
}

type authorize struct {
	db *gorm.DB
}

func newAuthorize(s *factory) Authorize {
	return &authorize{
		db: s.db.Model(&model.Authorize{}),
	}
}

// Create .
func (a *authorize) Create(ctx context.Context, data *model.Authorize) error {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	return a.db.WithContext(ctx).Create(data).Error
}

// Update .
func (a *authorize) Update(ctx context.Context, id int, data map[string]any) error {
	data["updated_at"] = time.Now()
	return a.db.WithContext(ctx).Where("id = ?", id).UpdateColumns(data).Error
}

// GetByProvider .
func (a *authorize) GetByProvider(ctx context.Context, provider, providerID string) (t *model.Authorize, err error) {
	err = a.db.WithContext(ctx).Where("provider = ? AND provider_id = ?", provider, providerID).Take(&t).Error
	return
}

// GetByUserID .
func (a *authorize) GetByUserID(ctx context.Context, userID int) (list []*model.Authorize, err error) {
	err = a.db.WithContext(ctx).Where("user_id = ?", userID).Find(&list).Error
	return
}

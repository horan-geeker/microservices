package model

import (
	"context"
	"gorm.io/gorm"
	"microservices/entity/model"
)

type File interface {
	Create(ctx context.Context, file *model.File) error
	Update(ctx context.Context, t *model.File) error
	GetList(ctx context.Context, userId, offset, limit int) (total int64, list []model.File, err error)
	GetByIdsAndUserId(ctx context.Context, ids []int, userId int) ([]model.File, error)
	GetByIdAndUserId(ctx context.Context, id, userId int) (t model.File, err error)
	DeleteByIdAndUserId(ctx context.Context, id, userId int) error
}

type file struct {
	db *gorm.DB
}

func newFile(db *gorm.DB) *file {
	return &file{
		db: db.Model(&model.File{}),
	}
}

func (f *file) GetList(ctx context.Context, userId, offset, limit int) (total int64, list []model.File, err error) {
	err = f.db.WithContext(ctx).Where("user_id", userId).Where("status", 1).Count(&total).Order("id desc").Offset(offset).Limit(limit).Find(&list).Error
	return
}

// Create .
func (f *file) Create(ctx context.Context, file *model.File) error {
	return f.db.WithContext(ctx).Create(file).Error
}

// Update .
func (f *file) Update(ctx context.Context, t *model.File) error {
	return f.db.WithContext(ctx).Updates(t).Error
}

// GetByIdsAndUserId .
func (f *file) GetByIdsAndUserId(ctx context.Context, ids []int, userId int) (list []model.File, err error) {
	err = f.db.WithContext(ctx).Where("id", ids).Where("status", 1).Where("user_id", userId).Find(&list).Error
	return
}

// GetByIdAndUserId .
func (f *file) GetByIdAndUserId(ctx context.Context, id, userId int) (t model.File, err error) {
	err = f.db.WithContext(ctx).Where("id", id).Where("status", 1).Where("user_id", userId).Take(&t).Error
	return
}

// DeleteByIdAndUserId .
func (f *file) DeleteByIdAndUserId(ctx context.Context, id, userId int) error {
	return f.db.WithContext(ctx).Where("id", id).Where("user_id", userId).Update("status", 0).Error
}

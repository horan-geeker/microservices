package model

import (
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"microservices/entity/model"
	"time"
)

// User defines the user storage interface.
type User interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, id uint64, data map[string]any) error
	GetByUid(ctx context.Context, id uint64) (*model.User, error)
	GetByName(ctx context.Context, name string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByPhone(ctx context.Context, phone string) (*model.User, error)
	//List(ctx context.Context) ([]model.User, error)
}

type user struct {
	db  *gorm.DB
	rdb *redis.Client
}

func newUser(s *factory) User {
	return &user{
		db: s.db.Model(&model.User{}),
	}
}

// Create .
func (u *user) Create(ctx context.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return u.db.WithContext(ctx).Create(&user).Error
}

// Update 更新用户表信息
func (u *user) Update(ctx context.Context, id uint64, data map[string]any) error {
	data["updated_at"] = time.Now()
	return u.db.WithContext(ctx).Where("id", id).UpdateColumns(data).Error
}

// GetByUid return an user by the user identifier.
func (u *user) GetByUid(ctx context.Context, id uint64) (t *model.User, err error) {
	err = u.db.WithContext(ctx).Where("id = ?", id).Take(&t).Error
	return
}

// GetByName .
func (u *user) GetByName(ctx context.Context, name string) (t *model.User, err error) {
	err = u.db.WithContext(ctx).Where("name = ?", name).Take(&t).Error
	return
}

// GetByEmail .
func (u *user) GetByEmail(ctx context.Context, email string) (t *model.User, err error) {
	err = u.db.WithContext(ctx).Where("email = ?", email).Take(&t).Error
	return
}

// GetByPhone .
func (u *user) GetByPhone(ctx context.Context, phone string) (t *model.User, err error) {
	err = u.db.WithContext(ctx).Where("phone = ?", phone).Take(&t).Error
	return
}

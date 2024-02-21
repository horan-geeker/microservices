package repository

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	consts2 "microservices/entity/consts"
	"microservices/entity/meta"
	"time"
)

// User defines the user storage interface.
type User interface {
	Create(ctx context.Context, user *meta.User) error
	Update(ctx context.Context, id uint64, data map[string]any) error
	GetByUid(ctx context.Context, id uint64) (*meta.User, error)
	GetByName(ctx context.Context, name string) (*meta.User, error)
	GetByEmail(ctx context.Context, email string) (*meta.User, error)
	GetByPhone(ctx context.Context, phone string) (*meta.User, error)
	//List(ctx context.Context) ([]model.User, error)
	SetToken(ctx context.Context, uid uint64, token string) error
	GetToken(ctx context.Context, uid uint64) (string, error)
	DeleteToken(ctx context.Context, uid uint64) error
}

type userImpl struct {
	db  *gorm.DB
	rdb *redis.Client
}

func newUsers(s *factoryImpl) User {
	return &userImpl{
		db:  s.db,
		rdb: s.rdb,
	}
}

// Create .
func (u *userImpl) Create(ctx context.Context, user *meta.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return u.db.WithContext(ctx).Create(&user).Error
}

// Update 更新用户表信息
func (u *userImpl) Update(ctx context.Context, id uint64, data map[string]any) error {
	data["updated_at"] = time.Now()
	return u.db.WithContext(ctx).Model(meta.User{}).Where("id", id).UpdateColumns(data).Error
}

// GetByUid return an user by the user identifier.
func (u *userImpl) GetByUid(ctx context.Context, id uint64) (*meta.User, error) {
	user := &meta.User{}
	err := u.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	return user, err
}

// GetByName .
func (u *userImpl) GetByName(ctx context.Context, name string) (*meta.User, error) {
	user := &meta.User{}
	err := u.db.WithContext(ctx).Where("name = ?", name).Take(&user).Error
	return user, err
}

// GetByEmail .
func (u *userImpl) GetByEmail(ctx context.Context, email string) (*meta.User, error) {
	user := &meta.User{}
	err := u.db.WithContext(ctx).Where("email = ?", email).Take(&user).Error
	return user, err
}

// GetByPhone .
func (u *userImpl) GetByPhone(ctx context.Context, phone string) (*meta.User, error) {
	user := &meta.User{}
	err := u.db.WithContext(ctx).Where("phone = ?", phone).Take(&user).Error
	return user, err
}

func (u *userImpl) SetToken(ctx context.Context, id uint64, token string) error {
	return u.rdb.Set(ctx, fmt.Sprintf(consts2.RedisUserTokenKey, id), token, consts2.UserTokenExpiredIn).Err()
}

func (u *userImpl) GetToken(ctx context.Context, id uint64) (string, error) {
	return u.rdb.Get(ctx, fmt.Sprintf(consts2.RedisUserTokenKey, id)).Result()
}

func (u *userImpl) DeleteToken(ctx context.Context, id uint64) error {
	return u.rdb.Del(ctx, fmt.Sprintf(consts2.RedisUserTokenKey, id)).Err()
}

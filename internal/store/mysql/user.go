package mysql

import (
	"context"
	"gorm.io/gorm"
	"microservices/internal/model"
	"microservices/internal/store"
	"time"
)

type users struct {
	db *gorm.DB
}

func newUsers(s *mysqlStore) store.UserStore {
	return &users{
		db: s.db,
	}
}

// Create .
func (u *users) Create(ctx context.Context, user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return u.db.WithContext(ctx).Create(&user).Error
}

// Update 更新用户表信息
func (u *users) Update(ctx context.Context, id uint64, data map[string]any) error {
	data["updated_at"] = time.Now()
	return u.db.WithContext(ctx).Model(model.User{}).Where("id", id).UpdateColumns(data).Error
}

// GetByUid return an user by the user identifier.
func (u *users) GetByUid(ctx context.Context, id uint64) (*model.User, error) {
	user := &model.User{}
	err := u.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	return user, err
}

// GetByName .
func (u *users) GetByName(ctx context.Context, name string) (*model.User, error) {
	user := &model.User{}
	err := u.db.WithContext(ctx).Where("name = ?", name).Take(&user).Error
	return user, err
}

// GetByEmail .
func (u *users) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := u.db.WithContext(ctx).Where("email = ?", email).Take(&user).Error
	return user, err
}

// GetByPhone .
func (u *users) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	user := &model.User{}
	err := u.db.WithContext(ctx).Where("phone = ?", phone).Take(&user).Error
	return user, err
}

package mysql

import (
	"context"
	"gorm.io/gorm"
	"microservices/internal/model"
	"microservices/internal/store"
)

type users struct {
	db *gorm.DB
}

func newUsers(s *mysqlStore) store.UserStore {
	return &users{
		db: s.db,
	}
}

// GetByUid return an user by the user identifier.
func (u *users) GetByUid(ctx context.Context, uid int) (*model.User, error) {
	user := &model.User{}
	err := u.db.WithContext(ctx).Where("id = ?", uid).Take(&user).Error
	return user, err
}

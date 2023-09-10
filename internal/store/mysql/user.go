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

func newUsers(ds *mysqlstore) store.UserStore {
	return &users{
		db: ds.db,
	}
}

// GetByUserName return an user by the user identifier.
func (u *users) GetByUserName(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := u.db.Where("name = ?", username).First(&user).Error
	return user, err
}

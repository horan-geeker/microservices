package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"microservices/internal/pkg/options"
	"microservices/internal/store"
)

type mysqlstore struct {
	db *gorm.DB
}

// Users .
func (s *mysqlstore) Users() store.UserStore {
	return newUsers(s)
}

// ConnectDB .
func ConnectDB(opts *options.MySQLOptions) store.Factory {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		opts.Username, opts.Password, opts.Host, opts.Port, opts.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return &mysqlstore{db: db}
}

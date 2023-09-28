package mysql

import (
	"gorm.io/gorm"
	"microservices/internal/store"
	"microservices/pkg/mysql"
	"sync"
)

// 定义 datastore
type datastore struct {
	db *gorm.DB
}

// Users .
func (s *datastore) Users() store.UserStore {
	return newUsers(s)
}

// 实例化
var (
	mysqlFactory store.Factory
	once         sync.Once
)

// GetMysqlInstance .
func GetMysqlInstance(opts *mysql.Options) store.Factory {
	once.Do(func() {
		db, err := mysql.NewMysql(opts)
		if err != nil {
			panic(err)
		}
		mysqlFactory = &datastore{
			db: db,
		}
	})
	return mysqlFactory
}

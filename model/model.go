package model

import (
	"gorm.io/gorm"
	"microservices/entity/config"
	"microservices/pkg/mysql"
)

// 实例化
var (
	factoryInstance Factory
	mysqlInstance   = GetMysqlInstance(config.NewMysqlOptions())
)

// Factory .
type Factory interface {
	User() User
}

// 定义 factory
type factory struct {
	db *gorm.DB
}

// User .
func (s *factory) User() User {
	return newUser(s)
}

// NewFactory .
func NewFactory() Factory {
	if factoryInstance == nil {
		factoryInstance = &factory{
			db: mysqlInstance,
		}
	}
	return factoryInstance
}

// GetMysqlInstance .
func GetMysqlInstance(opts *mysql.Options) *gorm.DB {
	mysqlInstance, err := mysql.NewMysql(opts)
	if err != nil {
		panic(err)
	}
	return mysqlInstance
}

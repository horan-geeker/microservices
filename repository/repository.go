package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"microservices/entity/config"
	"microservices/pkg/mysql"
	redis2 "microservices/pkg/redis"
	"sync"
)

// Factory .
type Factory interface {
	Users() User
	Auth() Auth
}

// 定义 factoryImpl
type factoryImpl struct {
	db  *gorm.DB
	rdb *redis.Client
}

// Users .
func (s *factoryImpl) Users() User {
	return newUsers(s)
}

// Auth .
func (s *factoryImpl) Auth() Auth {
	return newAuth(s)
}

// 实例化
var (
	factory       Factory
	db            *gorm.DB
	rdb           *redis.Client
	mysqlInitOnce sync.Once
	redisInitOnce sync.Once
)

// NewFactory .
func NewFactory() Factory {
	if factory == nil {
		factory = &factoryImpl{
			db:  GetMysqlInstance(config.NewMysqlOptions()),
			rdb: GetRedisInstance(config.NewRedisOptions()),
		}
	}
	return factory
}

// GetMysqlInstance .
func GetMysqlInstance(opts *mysql.Options) *gorm.DB {
	mysqlInitOnce.Do(func() {
		var err error
		db, err = mysql.NewMysql(opts)
		if err != nil {
			panic(err)
		}
	})
	return db
}

// GetRedisInstance .
func GetRedisInstance(opts *redis2.Options) *redis.Client {
	redisInitOnce.Do(func() {
		var err error
		rdb, err = redis2.NewRedis(opts)
		if err != nil {
			panic(err)
		}
	})
	return rdb
}

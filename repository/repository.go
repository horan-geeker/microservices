package repository

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"microservices/entity/options"
	"microservices/pkg/mysql"
	redis2 "microservices/pkg/redis"
	"sync"
)

// Factory .
type Factory interface {
	Users() User
	Auth() Auth
	Tencent() Tencent
	Aliyun() Aliyun
}

// 定义 factoryImpl
type factoryImpl struct {
	db         *gorm.DB
	rdb        *redis.Client
	tencentOpt *options.TencentOptions
	aliyunOpt  *options.AliyunOptions
}

// Users .
func (s *factoryImpl) Users() User {
	return newUsers(s)
}

// Auth .
func (s *factoryImpl) Auth() Auth {
	return newAuth(s)
}

func (s *factoryImpl) Tencent() Tencent {
	return newTencent(s.tencentOpt)
}

func (s *factoryImpl) Aliyun() Aliyun {
	return newAliyun(s.aliyunOpt)
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
			db:         GetMysqlInstance(options.NewMysqlOptions()),
			rdb:        GetRedisInstance(options.NewRedisOptions()),
			tencentOpt: options.NewTencentOptions(),
			aliyunOpt:  options.NewAliyunOptions(),
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

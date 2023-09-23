package router

import (
	"microservices/internal/pkg/options"
	"microservices/internal/store/mysql"
	"microservices/internal/store/redis"
	"microservices/pkg/meta"
)

var dataFactory = mysql.GetMysqlInstance(options.NewMySQLOptions())
var cacheFactory = redis.GetRedisInstance(options.NewRedisOptions())
var router = meta.GetEnginInstance()

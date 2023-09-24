package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"microservices/internal/middleware"
	"microservices/internal/pkg/options"
	"microservices/internal/store/mysql"
	"microservices/internal/store/redis"
	"microservices/pkg/meta"
)

var dataFactory = mysql.GetMysqlInstance(options.NewMySQLOptions())
var cacheFactory = redis.GetRedisInstance(options.NewRedisOptions())

// 注入全局中间件，注意压缩 response 的中间件顺序需要在 log response 之后压缩
var router = meta.GetEnginInstance(gzip.Gzip(gzip.DefaultCompression), middleware.RequestLogger, gin.Recovery())

package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"microservices/internal/pkg/options"
	"microservices/internal/store/redis"
	"microservices/pkg/util"
)

func Authenticate(c *gin.Context) {
	buf, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
	authorization := c.Request.Header.Get("Authorization")
	authClaims, err := util.ParseJWTToken(authorization)
	// 解析出错按照未登录返回
	if err != nil {

	}
	cache := redis.GetRedisInstance(options.NewRedisOptions())
	// 判断是否被注销
	_, err = cache.Users().GetToken(c.Request.Context(), authClaims.ID)
	if err != nil {

	}
}

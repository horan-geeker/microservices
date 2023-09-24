package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"microservices/internal/pkg/ecode"
	"microservices/internal/pkg/options"
	"microservices/internal/store/redis"
	"microservices/pkg/meta"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		buf, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
		authorization := c.Request.Header.Get("Authorization")
		if len(authorization) == 0 {
			meta.MakeErrorResponse(c, ecode.ErrTokenIsEmpty)
			c.Abort()
			return
		}
		authClaims, err := meta.ParseJWTToken(authorization)
		// 解析出错按照未登录返回
		if err != nil {
			meta.MakeErrorResponse(c, ecode.ErrTokenInvalid)
			c.Abort()
			return
		}
		cache := redis.GetRedisInstance(options.NewRedisOptions())
		// 判断是否被注销
		token, err := cache.Users().GetToken(c.Request.Context(), authClaims.Uid)
		if err != nil {
			meta.MakeErrorResponse(c, err)
			c.Abort()
			return
		}
		if token != authorization {
			meta.MakeErrorResponse(c, ecode.ErrTokenDiscard)
			c.Abort()
			return
		}
	}
}

package middleware

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	redis2 "github.com/redis/go-redis/v9"
	"io"
	"microservices/internal/pkg/ecode"
	"microservices/internal/pkg/jwt"
	"microservices/internal/pkg/options"
	"microservices/internal/store/redis"
	"microservices/pkg/app"
)

func Authenticate() gin.HandlerFunc {
	cache := redis.GetRedisInstance(options.NewRedisOptions())
	return func(c *gin.Context) {
		buf, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
		authorization := c.Request.Header.Get("Authorization")
		if len(authorization) == 0 {
			app.MakeErrorResponse(c, ecode.ErrTokenIsEmpty)
			c.Abort()
			return
		}
		authClaims, err := jwt.NewJwt(options.NewJwtOptions()).DecodeToken(authorization)
		// 解析出错按照未登录返回
		if err != nil {
			app.MakeErrorResponse(c, ecode.ErrTokenInvalid)
			c.Abort()
			return
		}
		// 判断是否被注销
		token, err := cache.Users().GetToken(c.Request.Context(), authClaims.Uid)
		if err == redis2.Nil {
			app.MakeErrorResponse(c, ecode.ErrTokenDiscard)
			c.Abort()
			return
		}
		if err != nil {
			app.MakeErrorResponse(c, err)
			c.Abort()
			return
		}
		if token != authorization {
			app.MakeErrorResponse(c, ecode.ErrTokenDiscard)
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "auth", authClaims))
	}
}

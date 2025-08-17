package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"microservices/cache"
	"microservices/entity/ecode"
	"microservices/entity/jwt"
	"microservices/pkg/app"
	"strconv"
	"time"
)

// ReqRateLimit 创建请求限频中间件
// limit: 限制次数
// seconds: 时间窗口（秒）
// useUID: true=根据用户ID限频，false=不使用用户ID限频
// useIP: true=根据IP限频，false=不使用IP限频
func ReqRateLimit(limit int, seconds int, useUID bool, useIP bool) gin.HandlerFunc {
	cacheFactory := cache.NewFactory()
	window := time.Duration(seconds) * time.Second
	keyPrefix := "req_rate_limit:"

	return func(c *gin.Context) {
		// 检查用户ID限频
		if useUID {
			if !checkUserRateLimit(c, cacheFactory, keyPrefix, limit, window, seconds) {
				return
			}
		}

		// 检查IP限频
		if useIP {
			if !checkIPRateLimit(c, cacheFactory, keyPrefix, limit, window, seconds) {
				return
			}
		}

		if !useUID && !useIP {
			// 检查全局限频
			if !checkGlobalRateLimit(c, cacheFactory, keyPrefix, limit, window, seconds) {
				return
			}
		}

		c.Next()
	}
}

// getUserRateLimitKey 获取基于用户ID的限频key
func getUserRateLimitKey(c *gin.Context, prefix string) string {
	// 尝试从context中获取auth信息
	authValue := c.Request.Context().Value("auth")
	if authValue == nil {
		return ""
	}

	authClaims, ok := authValue.(*jwt.AuthClaims)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%suser:%d", prefix, authClaims.Uid)
}

// getIPRateLimitKey 获取基于IP的限频key
func getIPRateLimitKey(c *gin.Context, prefix string) string {
	ip := c.ClientIP()
	return fmt.Sprintf("%sip:%s", prefix, ip)
}

// checkUserRateLimit 检查用户ID限频
func checkUserRateLimit(c *gin.Context, cacheFactory cache.Factory, keyPrefix string, limit int, window time.Duration, seconds int) bool {
	userKey := getUserRateLimitKey(c, keyPrefix)
	if userKey == "" {
		// 如果获取不到用户信息，跳过用户限频检查
		return true
	}

	allowed, count, err := cacheFactory.System().CheckRateLimit(c.Request.Context(), userKey, limit, window)
	if err != nil {
		// Redis错误时记录日志但不阻止请求
		c.Header("X-RateLimit-User-Error", "Redis error")
		return true
	}

	// 添加用户限频信息到响应头
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	c.Header("X-RateLimit-User-Limit", strconv.Itoa(limit))
	c.Header("X-RateLimit-User-Remaining", strconv.Itoa(remaining))
	c.Header("X-RateLimit-User-Window", strconv.Itoa(seconds))

	if !allowed {
		// 用户限频触发
		app.MakeErrorResponse(c, ecode.ErrTooManyRequests)
		c.Abort()
		return false
	}

	return true
}

// checkIPRateLimit 检查IP限频
func checkIPRateLimit(c *gin.Context, cacheFactory cache.Factory, keyPrefix string, limit int, window time.Duration, seconds int) bool {
	ipKey := getIPRateLimitKey(c, keyPrefix)

	allowed, count, err := cacheFactory.System().CheckRateLimit(c.Request.Context(), ipKey, limit, window)
	if err != nil {
		// Redis错误时记录日志但不阻止请求
		c.Header("X-RateLimit-IP-Error", "Redis error")
		return true
	}

	// 添加IP限频信息到响应头
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	c.Header("X-RateLimit-IP-Limit", strconv.Itoa(limit))
	c.Header("X-RateLimit-IP-Remaining", strconv.Itoa(remaining))
	c.Header("X-RateLimit-IP-Window", strconv.Itoa(seconds))

	if !allowed {
		// IP限频触发
		app.MakeErrorResponse(c, ecode.ErrTooManyRequests)
		c.Abort()
		return false
	}

	return true
}

// checkGlobalRateLimit 检查全局限频
func checkGlobalRateLimit(c *gin.Context, cacheFactory cache.Factory, keyPrefix string, limit int, window time.Duration, seconds int) bool {
	globalKey := fmt.Sprintf("%sglobal", keyPrefix)

	allowed, count, err := cacheFactory.System().CheckRateLimit(c.Request.Context(), globalKey, limit, window)
	if err != nil {
		// Redis错误时记录日志但不阻止请求
		c.Header("X-RateLimit-Global-Error", "Redis error")
		return true
	}

	// 添加全局限频信息到响应头
	remaining := limit - int(count)
	if remaining < 0 {
		remaining = 0
	}
	c.Header("X-RateLimit-Global-Limit", strconv.Itoa(limit))
	c.Header("X-RateLimit-Global-Remaining", strconv.Itoa(remaining))
	c.Header("X-RateLimit-Global-Window", strconv.Itoa(seconds))

	if !allowed {
		// 全局限频触发
		app.MakeErrorResponse(c, ecode.ErrTooManyRequests)
		c.Abort()
		return false
	}

	return true
}

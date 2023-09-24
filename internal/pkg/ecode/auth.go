package ecode

import (
	"errors"
	errors2 "microservices/pkg/ecode"
	"net/http"
)

var (
	// ErrTokenIsEmpty  用户未登录，token 为空
	ErrTokenIsEmpty = errors.New("用户未登录，token 为空")
	// ErrTokenInvalid  token无效
	ErrTokenInvalid = errors.New("用户未登录，token无效")
	// ErrTokenExpired  token已过期
	ErrTokenExpired = errors.New("用户未登录，token已过期")
	// ErrTokenNotExist  token不存在
	ErrTokenNotExist = errors.New("用户未登录，token不存在(已被注销)")
	// ErrTokenDiscard  token已注销
	ErrTokenDiscard = errors.New("用户未登录，token已注销")
	// ErrTokenInternalNotSet  token路由未设置
	ErrTokenInternalNotSet = errors.New("token路由未设置")
)

const (
	_ = iota + 2000
	// TokenIsEmpty  token为空
	TokenIsEmpty
	// TokenInvalid  token无效
	TokenInvalid
	// TokenExpired  token已过期
	TokenExpired
	// TokenNotExist  token不存在
	TokenNotExist
	// TokenDiscard  token已注销
	TokenDiscard
	// TokenInternalNotSet  token路由未设置
	TokenInternalNotSet
)

func init() {
	errors2.Register(ErrTokenIsEmpty, TokenIsEmpty, http.StatusUnauthorized)
	errors2.Register(ErrTokenInvalid, TokenInvalid, http.StatusUnauthorized)
	errors2.Register(ErrTokenExpired, TokenExpired, http.StatusUnauthorized)
	errors2.Register(ErrTokenNotExist, TokenNotExist, http.StatusUnauthorized)
	errors2.Register(ErrTokenDiscard, TokenDiscard, http.StatusUnauthorized)
	errors2.Register(ErrTokenInternalNotSet, TokenInternalNotSet, http.StatusInternalServerError)
}

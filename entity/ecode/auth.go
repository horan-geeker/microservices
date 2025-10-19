package ecode

import (
	"microservices/pkg/ecode"
	"net/http"
)

var (
	// ErrTokenIsEmpty 用户未登录，token 为空
	ErrTokenIsEmpty = &ecode.CustomError{Code: TokenIsEmpty, HttpStatus: http.StatusUnauthorized, Message: "用户未登录，token 为空"}
	// ErrTokenInvalid token无效
	ErrTokenInvalid = &ecode.CustomError{Code: TokenInvalid, HttpStatus: http.StatusUnauthorized, Message: "用户未登录，token无效"}
	// ErrTokenExpired token已过期
	ErrTokenExpired = &ecode.CustomError{Code: TokenExpired, HttpStatus: http.StatusUnauthorized, Message: "用户未登录，token已过期"}
	// ErrTokenNotExist token不存在
	ErrTokenNotExist = &ecode.CustomError{Code: TokenNotExist, HttpStatus: http.StatusUnauthorized, Message: "用户未登录，token不存在(已被注销)"}
	// ErrTokenDiscard token已注销
	ErrTokenDiscard = &ecode.CustomError{Code: TokenDiscard, HttpStatus: http.StatusUnauthorized, Message: "用户未登录，token已注销"}
	// ErrTokenInternalNotSet token路由未设置
	ErrTokenInternalNotSet = &ecode.CustomError{Code: TokenInternalNotSet, HttpStatus: http.StatusInternalServerError, Message: "token路由未设置"}
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

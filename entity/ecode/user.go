package ecode

import (
	"microservices/pkg/ecode"
	"net/http"
)

var (
	ErrInvalidPassword       = &ecode.CustomError{Code: UserPasswordError, HttpStatus: http.StatusOK, Message: "密码错误"}
	ErrUserPasswordDuplicate = &ecode.CustomError{Code: UserPasswordDuplicate, HttpStatus: http.StatusOK, Message: "设置新密码不能与旧密码相同"}
	ErrUserSmsCodeError      = &ecode.CustomError{Code: UserSmsCodeError, HttpStatus: http.StatusOK, Message: "短信验证码错误"}
	ErrUserEmailCodeError    = &ecode.CustomError{Code: UserEmailCodeError, HttpStatus: http.StatusOK, Message: "邮箱验证码错误"}
	ErrUserNameAlreadyExist  = &ecode.CustomError{Code: UserNameAlreadyExist, HttpStatus: http.StatusOK, Message: "用户名已存在"}
	ErrUserEmailAlreadyExist = &ecode.CustomError{Code: UserEmailAlreadyExist, HttpStatus: http.StatusOK, Message: "邮箱已存在"}
	ErrUserPhoneAlreadyExist = &ecode.CustomError{Code: UserPhoneAlreadyExist, HttpStatus: http.StatusOK, Message: "手机号已存在"}
	ErrAccountFrozen         = &ecode.CustomError{Code: AccountFrozen, HttpStatus: http.StatusForbidden, Message: "密码失败次数过多，账号已冻结，请两个小时后重试"}
)

const (
	_ = iota + 1000
	UserPasswordError
	UserPasswordDuplicate
	UserSmsCodeError
	UserEmailCodeError
	UserNameAlreadyExist
	UserEmailAlreadyExist
	UserPhoneAlreadyExist
	AccountFrozen
)

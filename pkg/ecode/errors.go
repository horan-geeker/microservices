package ecode

import (
	"errors"
	"net/http"
)

// CustomError 自定义错误类型，支持链式调用
type CustomError struct {
	Message    string
	Code       int
	HttpStatus int
}

// Error 实现error接口
func (e *CustomError) Error() string {
	return e.Message
}

// WithMessage 返回带有自定义消息的错误
func (e *CustomError) WithMessage(msg string) error {
	return &CustomError{
		Code:       e.Code,
		Message:    msg,
		HttpStatus: e.HttpStatus,
	}
}

// GetCode 获取错误码
func (e *CustomError) GetCode() int {
	return e.Code
}

// GetHttpStatus 获取HTTP状态码
func (e *CustomError) GetHttpStatus() int {
	return e.HttpStatus
}

// Is 实现 errors.Is 接口，基于错误码进行比较
func (e *CustomError) Is(target error) bool {
	var targetErr *CustomError
	if errors.As(target, &targetErr) {
		return e.Code == targetErr.Code
	}
	return false
}

var (
	ErrInternalServerError = &CustomError{Code: InternalServerErrorCode, Message: "系统内部错误", HttpStatus: http.StatusInternalServerError}
	ErrDataNotFound        = &CustomError{Code: DataNotFound, Message: "数据不存在", HttpStatus: http.StatusNotFound}
	ErrRouteParamInvalid   = &CustomError{Code: RouteParamInvalid, Message: "URL路由参数无效", HttpStatus: http.StatusBadRequest}
	ErrUserAuthFail        = &CustomError{Code: UserAuthFail, Message: "用户登录态校验失败", HttpStatus: http.StatusUnauthorized}
)

const (
	_ = iota
	InternalServerErrorCode
	DataNotFound
	RouteParamInvalid
	UserAuthFail
)

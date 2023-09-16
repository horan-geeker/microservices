package errors

import (
	"errors"
	"net/http"
)

var errMapErrCode = make(map[error]int)
var errMapHttpStatus = make(map[error]int)
var errCollect []error

// Register .
func Register(err error, errCode int, httpStatus int) {
	errCollect = append(errCollect, err)
	errMapErrCode[err] = errCode
	errMapHttpStatus[err] = httpStatus
}

// GetCollectErr .
func GetCollectErr() []error {
	return errCollect
}

// GetErrCodeByErr .
func GetErrCodeByErr(err error) int {
	if code, ok := errMapErrCode[err]; ok {
		return code
	}
	return InternalServerErrorCode
}

// GetHttpStatusByErr .
func GetHttpStatusByErr(err error) int {
	if status, ok := errMapHttpStatus[err]; ok {
		return status
	}
	return http.StatusInternalServerError
}

var (
	ErrInternalServerError = errors.New("系统内部错误")
	ErrDataNotFound        = errors.New("数据不存在")
)

const (
	_ = iota
	InternalServerErrorCode
	DataNotFound
)

func init() {
	Register(ErrInternalServerError, InternalServerErrorCode, http.StatusInternalServerError)
	Register(ErrDataNotFound, DataNotFound, http.StatusNotFound)
}

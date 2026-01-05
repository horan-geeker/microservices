package ecode

import (
	"microservices/pkg/ecode"
	"net/http"
)

var (
	// ErrGenerationCreditNotEnough 用户积分不足
	ErrGenerationCreditNotEnough = &ecode.CustomError{Code: GenerationCreditNotEnough, HttpStatus: http.StatusOK, Message: "user credit not enough"}
)

const (
	_ = iota + 5000
	GenerationCreditNotEnough
)

package util

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"microservices/internal/pkg/options"
	"microservices/pkg/meta"
)

// MD5 md5 加密
func MD5(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func ParseJWTToken(tokenString string) (*meta.AuthClaims, error) {
	// 解析token
	opts := options.NewJwtOptions()
	token, err := jwt.ParseWithClaims(tokenString, &meta.AuthClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(opts.Key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*meta.AuthClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

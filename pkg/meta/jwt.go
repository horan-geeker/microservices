package meta

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"microservices/internal/pkg/options"
)

type AuthClaims struct {
	Uid uint64 `json:"uid"`
	jwt.StandardClaims
}

func ParseJWTToken(tokenString string) (*AuthClaims, error) {
	// 解析token
	opts := options.NewJwtOptions()
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(opts.Key), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("token is invalid")
}

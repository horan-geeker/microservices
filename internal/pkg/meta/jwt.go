package meta

import "github.com/golang-jwt/jwt"

type AuthClaims struct {
	Uid uint64 `json:"uid"`
	jwt.StandardClaims
}

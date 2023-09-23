package meta

import "github.com/golang-jwt/jwt"

type AuthClaims struct {
	ID uint64 `json:"id"`
	jwt.StandardClaims
}

package auth

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt.Claims
	ID int
}

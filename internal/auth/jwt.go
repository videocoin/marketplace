package auth

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	jwt.StandardClaims
	IsActive bool   `json:"is_active"`
	Address  string `json:"address"`
}

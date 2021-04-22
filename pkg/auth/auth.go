package auth

import (
	"context"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Claims struct {
	jwt.StandardClaims
	IsActive bool   `json:"is_active"`
	Address  string `json:"address"`
}

func AuthFromContext(ctx context.Context) (context.Context, error) {
	tokenStr, err := grpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, ErrInvalidToken
	}

	secret, ok := SecretKeyFromContext(ctx)
	if !ok {
		return ctx, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return ctx, err
	}

	if !token.Valid {
		return ctx, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return ctx, ErrInvalidToken
	}

	ctx = NewContextWithJWTClaims(ctx, claims)
	return ctx, nil
}

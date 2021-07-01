package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

type key int

const (
	secretKey     key = 0
	jwtSubjectKey key = 1
	jwtClaimsKey  key = 2
)

func NewContextWithSecretKey(ctx context.Context, secret string) context.Context {
	return context.WithValue(ctx, secretKey, secret)
}

func SecretKeyFromContext(ctx context.Context) (string, bool) {
	secret, ok := ctx.Value(secretKey).(string)
	return secret, ok
}

func NewContextWithJWTSubject(ctx context.Context, sub string) context.Context {
	return context.WithValue(ctx, jwtSubjectKey, sub)
}

func JWTSubjectFromContext(ctx context.Context) (string, bool) {
	sub, ok := ctx.Value(jwtSubjectKey).(string)
	return sub, ok
}

func NewContextWithJWTClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, jwtClaimsKey, claims)
}

func JWTClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(jwtClaimsKey).(*Claims)
	return claims, ok
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


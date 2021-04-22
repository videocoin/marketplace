package auth

import (
	"context"
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

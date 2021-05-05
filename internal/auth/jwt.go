package auth

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/videocoin/marketplace/internal/model"
	"strconv"
	"time"
)

type Claims struct {
	jwt.StandardClaims
	IsActive bool   `json:"is_active"`
	Address  string `json:"address"`
}

func CreateAuthToken(ctx context.Context, secret string, account *model.Account) (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.FormatInt(account.ID, 10),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
		IsActive: account.IsActive,
		Address: account.Address,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
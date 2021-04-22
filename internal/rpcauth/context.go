package rpcauth

import (
	"context"
	"github.com/videocoin/marketplace/internal/model"
)

type key int

const (
	accountKey key = 1
)

func NewContextWithAccount(ctx context.Context, user *model.Account) context.Context {
	return context.WithValue(ctx, accountKey, user)
}

func AccountFromContext(ctx context.Context) (*model.Account, bool) {
	user, ok := ctx.Value(accountKey).(*model.Account)
	return user, ok
}

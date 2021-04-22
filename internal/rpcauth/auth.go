package rpcauth

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/videocoin/marketplace/api/rpc"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/pkg/auth"
	"strings"
)

var (
	NoAuthSuffixMethods = []string{
		".Health",
		".AccountsService/GetNonce",
		".AccountsService/Register",
		".AccountsService/Auth",
		".MarketplaceService/GetArts",
		".MarketplaceService/GetArt",
		".MarketplaceService/GetCreators",
		".MarketplaceService/GetCreator",
		".MarketplaceService/GetArtsByCreator",
		".MarketplaceService/GetSpotlightFeaturedArts",
		".MarketplaceService/GetSpotlightLiveArts",
		".MarketplaceService/GetSpotlightFeaturedCreators",
	}
)

func Auth(
	ctx context.Context,
	fullMethodName string,
	secret string,
	ds *datastore.Datastore,
) (context.Context, error) {
	for _, suffix := range NoAuthSuffixMethods {
		if strings.HasSuffix(fullMethodName, suffix) {
			return ctx, nil
		}
	}

	ctx = auth.NewContextWithSecretKey(ctx, secret)

	ctx, err := auth.AuthFromContext(ctx)
	if err != nil {
		ctxlogrus.Extract(ctx).WithError(err).Warning("failed to auth from context")
		return ctx, rpc.ErrRpcUnauthenticated
	}

	claims, _ := auth.JWTClaimsFromContext(ctx)
	if claims.Subject == "" {
		ctxlogrus.Extract(ctx).WithError(err).Warning("failed to get subject from context")
		return ctx, rpc.ErrRpcUnauthenticated
	}

	account, err := ds.Accounts.GetByAddress(ctx, claims.Address)
	if err != nil {
		if rpc.IsNotFoundError(err) {
			ctxlogrus.Extract(ctx).WithError(err).Error("account not found")
			return nil, rpc.ErrRpcUnauthenticated
		}
		ctxlogrus.Extract(ctx).WithError(err).Error("failed to get account")
		return nil, rpc.ErrRpcUnauthenticated
	}

	return NewContextWithAccount(ctx, account), nil
}

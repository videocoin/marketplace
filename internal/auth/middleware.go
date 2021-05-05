package auth

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"net/http"
)

func JWTAuth(logger *logrus.Entry, ds *datastore.Datastore, secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := NewContextWithSecretKey(context.Background(), secret)
			ah := c.Request().Header.Get("Authorization")
			md := metautils.NiceMD{"authorization": []string{ah}}

			ctx, err := AuthFromContext(md.ToIncoming(ctx))
			if err != nil {
				logger.WithError(err).Warning("failed to parse jwt token")

				return &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  "invalid or expired jwt",
					Internal: err,
				}
			}

			claims, _ := JWTClaimsFromContext(ctx)
			account, err := ds.Accounts.GetByAddress(ctx, claims.Address)
			if err != nil {
				logger.WithError(err).Error("failed to get account")

				return &echo.HTTPError{
					Code:     http.StatusUnauthorized,
					Message:  "invalid or expired jwt",
					Internal: err,
				}
			}

			c.Set("account", account)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

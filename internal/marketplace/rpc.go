package marketplace

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/rpcauth"
	"github.com/videocoin/marketplace/pkg/grpcutil"
)

type Server struct {
	logger     *logrus.Entry
	validator  *grpcutil.RequestValidator
	authSecret string
	ds         *datastore.Datastore
}

func NewServer(ctx context.Context, opts ...ServerOption) *Server {
	s := &Server{
		logger: ctxlogrus.Extract(ctx).WithField("system", "marketplace"),
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

func (s *Server) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return rpcauth.Auth(ctx, fullMethodName, s.authSecret, s.ds)
}

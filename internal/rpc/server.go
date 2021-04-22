package rpc

import (
	"context"
	"github.com/sirupsen/logrus"
	accountsv1 "github.com/videocoin/marketplace/api/v1/accounts"
	marketplacev1 "github.com/videocoin/marketplace/api/v1/marketplace"
	"github.com/videocoin/marketplace/internal/accounts"
	"github.com/videocoin/marketplace/internal/marketplace"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/pkg/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server struct {
	logger     *logrus.Entry
	addr       string
	lis        net.Listener
	grpcSrv    *grpc.Server
	grpcOpts   []grpc.ServerOption
	authSecret string
	ds         *datastore.Datastore
}

func NewServer(ctx context.Context, opts ...Option) (*Server, error) {
	srv := new(Server)

	for _, o := range opts {
		if err := o(srv); err != nil {
			return nil, err
		}
	}

	srv.grpcSrv = grpc.NewServer(srv.grpcOpts...)

	lis, err := net.Listen("tcp", srv.addr)
	if err != nil {
		return nil, err
	}
	srv.lis = lis

	validator, err := grpcutil.NewRequestValidator()
	if err != nil {
		return nil, err
	}

	healthv1.RegisterHealthServer(srv.grpcSrv, health.NewServer())
	accountsv1.RegisterAccountsServiceServer(
		srv.grpcSrv,
		accounts.NewServer(
			ctx,
			accounts.WithValidator(validator),
			accounts.WithDatastore(srv.ds),
			accounts.WithAuthSecret(srv.authSecret),
		),
	)
	marketplacev1.RegisterMarketplaceServiceServer(
		srv.grpcSrv,
		marketplace.NewServer(
			ctx,
			marketplace.WithValidator(validator),
			marketplace.WithDatastore(srv.ds),
			marketplace.WithAuthSecret(srv.authSecret),
		),
	)

	reflection.Register(srv.grpcSrv)

	return srv, nil
}

func (s *Server) Start() error {
	s.logger.WithField("addr", s.addr).Info("starting rpc server")
	return s.grpcSrv.Serve(s.lis)
}

func (s *Server) Stop() error {
	s.logger.Infof("stopping rpc server")
	s.grpcSrv.GracefulStop()
	return nil
}

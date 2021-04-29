package rpc

import (
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/storage"
	"google.golang.org/grpc"
)

type Option func(*Server) error

func WithAddr(addr string) Option {
	return func(s *Server) error {
		s.addr = addr
		return nil
	}
}

func WithLogger(logger *logrus.Entry) Option {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}

func WithGRPCServerOpts(opts []grpc.ServerOption) Option {
	return func(s *Server) error {
		s.grpcOpts = opts
		return nil
	}
}

func WithDatastore(ds *datastore.Datastore) Option {
	return func(s *Server) error {
		s.ds = ds
		return nil
	}
}

func WithAuthSecret(secret string) Option {
	return func(s *Server) error {
		s.authSecret = secret
		return nil
	}
}

func WithStorage(storage *storage.Storage) Option {
	return func(s *Server) error {
		s.storage = storage
		return nil
	}
}
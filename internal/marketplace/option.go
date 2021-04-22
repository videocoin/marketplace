package marketplace

import (
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/pkg/grpcutil"
)

type ServerOption func(*Server)

func WithValidator(v *grpcutil.RequestValidator) ServerOption {
	return func(s *Server) {
		s.validator = v
	}
}

func WithDatastore(ds *datastore.Datastore) ServerOption {
	return func(s *Server) {
		s.ds = ds
	}
}

func WithAuthSecret(secret string) ServerOption {
	return func(s *Server) {
		s.authSecret = secret
	}
}

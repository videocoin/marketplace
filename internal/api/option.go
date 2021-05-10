package api

import (
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/minter"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/internal/mediaconverter"
)

type ServerOption func(service *Server) error

func WithLogger(logger *logrus.Entry) ServerOption {
	return func(s *Server) error {
		s.logger = logger
		return nil
	}
}

func WithAddr(addr string) ServerOption {
	return func(s *Server) error {
		s.addr = addr
		return nil
	}
}

func WithAuthSecret(secret string) ServerOption {
	return func(s *Server) error {
		s.authSecret = secret
		return nil
	}
}

func WithDatastore(ds *datastore.Datastore) ServerOption {
	return func(s *Server) error {
		s.ds = ds
		return nil
	}
}

func WithGCPBucket(bucket string) ServerOption {
	return func(s *Server) error {
		s.gcpBucket = bucket
		return nil
	}
}

func WithMediaConverter(mc *mediaconverter.MediaConverter) ServerOption {
	return func(s *Server) error {
		s.mc = mc
		return nil
	}
}

func WithStorage(storage *storage.Storage) ServerOption {
	return func(s *Server) error {
		s.storage = storage
		return nil
	}
}

func WithMinter(m *minter.Minter) ServerOption {
	return func(s *Server) error {
		s.minter = m
		return nil
	}
}
package assets

import (
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/storage"
	"github.com/videocoin/marketplace/internal/mediaconverter"
)

type Option func(service *AssetsService) error

func WithLogger(logger *logrus.Entry) Option {
	return func(s *AssetsService) error {
		s.logger = logger
		return nil
	}
}

func WithAuthSecret(secret string) Option {
	return func(s *AssetsService) error {
		s.authSecret = secret
		return nil
	}
}

func WithDatastore(ds *datastore.Datastore) Option {
	return func(s *AssetsService) error {
		s.ds = ds
		return nil
	}
}

func WithGCPBucket(bucket string) Option {
	return func(s *AssetsService) error {
		s.gcpBucket = bucket
		return nil
	}
}

func WithMediaConverter(mc *mediaconverter.MediaConverter) Option {
	return func(s *AssetsService) error {
		s.mc = mc
		return nil
	}
}

func WithStorage(storage *storage.Storage) Option {
	return func(s *AssetsService) error {
		s.storage = storage
		return nil
	}
}

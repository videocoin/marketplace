package assets

import (
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
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

func WithIPFSGateway(addr string) Option {
	return func(s *AssetsService) error {
		s.ipfsGw = addr
		return nil
	}
}

func WithBucket(bucket string) Option {
	return func(s *AssetsService) error {
		s.bucket = bucket
		return nil
	}
}

func WithMediaConverter(mc *mediaconverter.MediaConverter) Option {
	return func(s *AssetsService) error {
		s.mc = mc
		return nil
	}
}

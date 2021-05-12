package mediaconverter

import (
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/storage"
)

type Option func(*MediaConverter) error

func WithLogger(logger *logrus.Entry) Option {
	return func(mc *MediaConverter) error {
		mc.logger = logger
		return nil
	}
}

func WithDatastore(ds *datastore.Datastore) Option {
	return func(mc *MediaConverter) error {
		mc.ds = ds
		return nil
	}
}

func WithGCPConfig(config *GCPConfig) Option {
	return func(mc *MediaConverter) error {
		mc.gcpConfig = config
		return nil
	}
}

func WithStorage(s *storage.Storage) Option {
	return func(mc *MediaConverter) error {
		mc.storage = s
		return nil
	}
}

func WithTranscoding() Option {
	return func(mc *MediaConverter) error {
		mc.enableTranscoding = true
		return nil
	}
}

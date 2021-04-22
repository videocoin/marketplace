package mediaconverter

import (
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
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

func WithIPFSGateway(addr string) Option {
	return func(mc *MediaConverter) error {
		mc.ipfsGw = addr
		return nil
	}
}

func WithGCPConfig(config *GCPConfig) Option {
	return func(mc *MediaConverter) error {
		mc.gcpConfig = config
		return nil
	}
}

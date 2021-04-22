package gateway

import (
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/assets"
)

type Option func(*Gateway) error

func WithLogger(logger *logrus.Entry) Option {
	return func(gw *Gateway) error {
		gw.logger = logger
		return nil
	}
}

func WithAddr(addr string) Option {
	return func(gw *Gateway) error {
		gw.addr = addr
		return nil
	}
}

func WithBackendAddr(addr string) Option {
	return func(gw *Gateway) error {
		gw.backendAddr = addr
		return nil
	}
}

func WithAssetsService(svc *assets.AssetsService) Option {
	return func(gw *Gateway) error {
		gw.assetsSvc = svc
		return nil
	}
}

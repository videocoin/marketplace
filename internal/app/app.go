package app

import (
	"context"
	"github.com/videocoin/marketplace/internal/storage"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/assets"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/gateway"
	"github.com/videocoin/marketplace/internal/mediaconverter"
	"github.com/videocoin/marketplace/internal/rpc"
	"github.com/videocoin/marketplace/pkg/grpcutil"
)

type App struct {
	cfg    *Config
	logger *logrus.Entry
	stop   chan bool
	ds     *datastore.Datastore
	rpcSrv *rpc.Server
	gw     *gateway.Gateway
	mc     *mediaconverter.MediaConverter
}

func NewApp(ctx context.Context, cfg *Config) (*App, error) {
	logger := ctxlogrus.Extract(ctx)

	dsCtx := ctxlogrus.ToContext(ctx, logger.WithField("system", "datastore"))
	ds, err := datastore.NewDatastore(dsCtx, cfg.DBURI)
	if err != nil {
		return nil, err
	}

	storageCli, err := storage.NewStorage(storage.WithConfig(&storage.TextileConfig{
		AuthKey:       cfg.TextileAuthKey,
		AuthSecret:    cfg.TextileAuthSecret,
		ThreadID:      cfg.TextileThreadID,
		BucketRootKey: cfg.TextileBucketRootKey,
	}))
	if err != nil {
		return nil, err
	}

	mc, err := mediaconverter.NewMediaConverter(
		ctx,
		mediaconverter.WithLogger(logger.WithField("system", "mediaconverter")),
		mediaconverter.WithDatastore(ds),
		mediaconverter.WithStorage(storageCli),
		mediaconverter.WithGCPConfig(&mediaconverter.GCPConfig{
			Bucket:             cfg.GCPBucket,
			Project:            cfg.GCPProject,
			Region:             cfg.GCPRegion,
			PubSubTopic:        cfg.GCPPubSubTopic,
			PubSubSubscription: cfg.GCPPubSubSubscription,
		}),
	)
	if err != nil {
		return nil, err
	}

	rpcSrv, err := rpc.NewServer(
		ctx,
		rpc.WithAddr(cfg.RPCAddr),
		rpc.WithLogger(logger.WithField("system", "rpc")),
		rpc.WithGRPCServerOpts(grpcutil.DefaultServerOpts(logger)),
		rpc.WithAuthSecret(cfg.AuthSecret),
		rpc.WithDatastore(ds),
		rpc.WithStorage(storageCli),
	)
	if err != nil {
		return nil, err
	}

	assetsSvc, err := assets.NewAssetsService(
		ctx,
		assets.WithLogger(logger.WithField("system", "assets")),
		assets.WithAuthSecret(cfg.AuthSecret),
		assets.WithDatastore(ds),
		assets.WithStorage(storageCli),
		assets.WithGCPBucket(cfg.GCPBucket),
		assets.WithMediaConverter(mc),
	)
	if err != nil {
		return nil, err
	}

	gw, err := gateway.NewGateway(
		ctx,
		gateway.WithLogger(logger.WithField("system", "gateway")),
		gateway.WithAddr(cfg.GWAddr),
		gateway.WithBackendAddr(cfg.RPCAddr),
		gateway.WithAssetsService(assetsSvc),
	)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:    cfg,
		logger: ctxlogrus.Extract(ctx),
		stop:   make(chan bool, 1),
		ds:     ds,
		rpcSrv: rpcSrv,
		gw:     gw,
		mc:     mc,
	}, nil
}

func (s *App) Start(errCh chan error) {
	go func() {
		err := s.rpcSrv.Start()
		if err != nil {
			errCh <- err
		}
	}()

	go func() {
		s.gw.Start(errCh)
	}()

	go func() {
		s.mc.Start(errCh)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			s.logger.WithError(err).Error("failed to start")
		}
	}
}

func (s *App) Stop() error {
	err := s.rpcSrv.Stop()
	if err != nil {
		s.logger.WithError(err).Error("failed to stop rpc server")
	}

	err = s.gw.Stop()
	if err != nil {
		s.logger.WithError(err).Error("failed to stop gateway server")
	}

	err = s.mc.Stop()
	if err != nil {
		s.logger.WithError(err).Error("failed to stop media converter")
	}

	s.stop <- true
	return nil
}

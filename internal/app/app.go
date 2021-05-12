package app

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/api"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/mediaconverter"
	minter "github.com/videocoin/marketplace/internal/minter"
	"github.com/videocoin/marketplace/internal/storage"
)

type App struct {
	cfg    *Config
	logger *logrus.Entry
	stop   chan bool
	ds     *datastore.Datastore
	mc     *mediaconverter.MediaConverter
	api    *api.Server
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

	mcOpts := []mediaconverter.Option{
		mediaconverter.WithLogger(logger.WithField("system", "mediaconverter")),
		mediaconverter.WithDatastore(ds),
		mediaconverter.WithStorage(storageCli),
	}
	if cfg.EnableTranscoding {
		mcOpts = append(mcOpts, mediaconverter.WithGCPConfig(&mediaconverter.GCPConfig{
			Bucket:             cfg.GCPBucket,
			Project:            cfg.GCPProject,
			Region:             cfg.GCPRegion,
			PubSubTopic:        cfg.GCPPubSubTopic,
			PubSubSubscription: cfg.GCPPubSubSubscription,
		}), mediaconverter.WithTranscoding())
	}
	mc, err := mediaconverter.NewMediaConverter(
		ctx,
		mcOpts...,
	)
	if err != nil {
		return nil, err
	}

	m, err := minter.NewMinter(
		cfg.BlockchainURL,
		cfg.ERC1155ContractAddress,
		cfg.ERC1155ContractKeyFile,
		cfg.ERC1155ContractKeyPass,
	)
	if err != nil {
		return nil, err
	}

	apiSrv, err := api.NewServer(
		ctx,
		api.WithAddr(cfg.Addr),
		api.WithLogger(logger.WithField("system", "assets")),
		api.WithAuthSecret(cfg.AuthSecret),
		api.WithDatastore(ds),
		api.WithStorage(storageCli),
		api.WithGCPBucket(cfg.GCPBucket),
		api.WithMediaConverter(mc),
		api.WithMinter(m),
	)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:    cfg,
		logger: ctxlogrus.Extract(ctx),
		stop:   make(chan bool, 1),
		ds:     ds,
		mc:     mc,
		api:    apiSrv,
	}, nil
}

func (s *App) Start(errCh chan error) {
	go func() {
		s.api.Start(errCh)
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
	err := s.api.Stop()
	if err != nil {
		s.logger.WithError(err).Error("failed to stop api server")
	}

	err = s.mc.Stop()
	if err != nil {
		s.logger.WithError(err).Error("failed to stop media converter")
	}

	s.stop <- true
	return nil
}

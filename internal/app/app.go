package app

import (
	"context"
	"github.com/videocoin/marketplace/internal/mediaprocessor"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/api"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/listener"
	"github.com/videocoin/marketplace/internal/minter"
	"github.com/videocoin/marketplace/internal/orderbook"
	"github.com/videocoin/marketplace/internal/storage"
)

type App struct {
	cfg    *Config
	logger *logrus.Entry
	stop   chan bool
	ds     *datastore.Datastore
	mp     *mediaprocessor.MediaProcessor
	api    *api.Server
	el     *listener.ExchangeListener
}

func NewApp(ctx context.Context, cfg *Config) (*App, error) {
	logger := ctxlogrus.Extract(ctx)

	dsCtx := ctxlogrus.ToContext(ctx, logger.WithField("system", "datastore"))
	ds, err := datastore.NewDatastore(dsCtx, cfg.DBURI)
	if err != nil {
		return nil, err
	}

	storageOpts := make([]storage.Option, 0)
	if cfg.StorageBackend == storage.NftStorage {
		storageOpts = append(storageOpts, storage.WithNftStorage(&storage.NftStorageConfig{
			ApiKey: cfg.NftStorageApiKey,
		}))
		logger.WithField("storage", "nftstorage").Info("initializing storage backend")
	} else {
		storageOpts = append(storageOpts, storage.WithTextile(&storage.TextileConfig{
			AuthKey:       cfg.TextileAuthKey,
			AuthSecret:    cfg.TextileAuthSecret,
			ThreadID:      cfg.TextileThreadID,
			BucketRootKey: cfg.TextileBucketRootKey,
		}))
		logger.WithField("storage", "textile").Info("initializing storage backend")
	}
	storageCli, err := storage.NewStorage(storageOpts...)
	if err != nil {
		return nil, err
	}

	mpOpts := []mediaprocessor.Option{
		mediaprocessor.WithLogger(logger.WithField("system", "mediaprocessor")),
		mediaprocessor.WithDatastore(ds),
		mediaprocessor.WithStorage(storageCli),
	}

	mc, err := mediaprocessor.NewMediaProcessor(ctx, mpOpts...)
	if err != nil {
		return nil, err
	}

	m, err := minter.NewMinter(
		cfg.BlockchainURL,
		cfg.ERC721ContractAddress,
		cfg.ERC721ContractKeyFile,
		cfg.ERC721ContractKeyPass,
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
		api.WithMediaConverter(mc),
		api.WithMinter(m),
	)
	if err != nil {
		return nil, err
	}

	ob, err := orderbook.NewOderBook(
		ctx,
		orderbook.WithMinter(m),
		orderbook.WithDatastore(ds),
		orderbook.WithMediaProcessor(mc),
		orderbook.WithStorage(storageCli),
	)
	if err != nil {
		return nil, err
	}

	el, err := listener.NewExchangeListener(
		ctx,
		listener.WithBlockchainURL(cfg.BlockchainURL),
		listener.WithScanFrom(cfg.BlockchainScanFrom),
		listener.WithContractAddress(cfg.ERC721AuctionContractAddress),
		listener.WithDatastore(ds),
		listener.WithOrderbook(ob),
	)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:    cfg,
		logger: ctxlogrus.Extract(ctx),
		stop:   make(chan bool, 1),
		ds:     ds,
		mp:     mc,
		api:    apiSrv,
		el:     el,
	}, nil
}

func (s *App) Start(errCh chan error) {
	go func() {
		s.api.Start(errCh)
	}()

	go func() {
		s.el.Start(errCh)
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

	err = s.el.Stop()
	if err != nil {
		s.logger.WithError(err).Error("failed to stop exchange listener")
	}

	s.stop <- true
	return nil
}

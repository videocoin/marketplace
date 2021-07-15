package orderbook

import (
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/mediaprocessor"
	"github.com/videocoin/marketplace/internal/minter"
	"github.com/videocoin/marketplace/internal/storage"
)

type Option func(l *OrderBook) error

func WithLogger(logger *logrus.Entry) Option {
	return func(book *OrderBook) error {
		book.logger = logger
		return nil
	}
}

func WithDatastore(ds *datastore.Datastore) Option {
	return func(book *OrderBook) error {
		book.ds = ds
		return nil
	}
}

func WithMediaProcessor(mp *mediaprocessor.MediaProcessor) Option {
	return func(book *OrderBook) error {
		book.mp = mp
		return nil
	}
}

func WithStorage(s *storage.Storage) Option {
	return func(book *OrderBook) error {
		book.storage = s
		return nil
	}
}

func WithMinter(m *minter.Minter) Option {
	return func(book *OrderBook) error {
		book.minter = m
		return nil
	}
}

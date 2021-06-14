package listener

import (
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
)

type ExchangeListenerOption func(l *ExchangeListener) error

func WithLogger(logger *logrus.Entry) ExchangeListenerOption {
	return func(l *ExchangeListener) error {
		l.logger = logger
		return nil
	}
}

func WithDatastore(ds *datastore.Datastore) ExchangeListenerOption {
	return func(l *ExchangeListener) error {
		l.ds = ds
		return nil
	}
}

func WithBlockchainURL(u string) ExchangeListenerOption {
	return func(l *ExchangeListener) error {
		l.url = u
		return nil
	}
}

func WithContractAddress(ca string) ExchangeListenerOption {
	return func(l *ExchangeListener) error {
		l.ca = ca
		return nil
	}
}

func WithLogStep(step uint64) ExchangeListenerOption {
	return func(l *ExchangeListener) error {
		l.logStep = step
		return nil
	}
}

func WithScanFrom(scanFrom uint64) ExchangeListenerOption {
	return func(l *ExchangeListener) error {
		l.scanFrom = scanFrom
		return nil
	}
}

package listener

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/videocoin/marketplace/internal/model"
	"github.com/videocoin/marketplace/internal/wyvern"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/marketplace/internal/datastore"
	"github.com/videocoin/marketplace/internal/orderbook"
)

type ExchangeListener struct {
	logger    *logrus.Entry
	ds        *datastore.Datastore
	orderbook *orderbook.OrderBook
	url       string
	ca        string
	logStep   uint64
	scanFrom  uint64
	cli       *ethclient.Client
	re        *EventReader
	chainID   string
	t         *time.Ticker
}

func NewExchangeListener(ctx context.Context, opts ...ExchangeListenerOption) (*ExchangeListener, error) {
	l := &ExchangeListener{
		logger:   ctxlogrus.Extract(ctx).WithField("system", "exchange-listener"),
		logStep:  1000,
		scanFrom: 0,
		t:        time.NewTicker(time.Second * 5),
	}

	for _, o := range opts {
		if err := o(l); err != nil {
			return nil, err
		}
	}

	chainIDHash := md5.Sum([]byte(fmt.Sprintf("%s#%s", l.url, l.ca)))
	l.chainID = hex.EncodeToString(chainIDHash[:])

	cli, err := ethclient.Dial(l.url)
	if err != nil {
		return nil, err
	}

	l.cli = cli

	erLogger := l.logger.WithField("system", "event-reader")
	re, err := NewEventReader(cli, l.ca, erLogger)
	if err != nil {
		return nil, err
	}

	l.re = re

	_, err = l.ds.ChainMeta.GetLastHeight(ctx, l.chainID)
	if err == datastore.ErrChainMetaNotFound {
		initErr := l.ds.ChainMeta.Init(ctx, l.chainID)
		if initErr != nil {
			return nil, err
		}
	}

	return l, nil
}

func (listener *ExchangeListener) headNumber(ctx context.Context) (uint64, error) {
	header, err := listener.cli.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

func (listener *ExchangeListener) waitEvents(ctx context.Context) error {
	knownHeight, err := listener.ds.ChainMeta.GetLastHeight(ctx, listener.chainID)
	if err != nil {
		return err
	}

	number, err := listener.headNumber(ctx)
	if err != nil {
		return err
	}

	var start uint64
	if knownHeight < listener.scanFrom {
		start = listener.scanFrom
	} else {
		start = knownHeight + 1
	}

	if start > number {
		start = number
	}

	end := start + listener.logStep
	if end > number {
		end = number
	}

	listener.logger.
		WithField("block_start", start).
		WithField("block_end", end).
		Info("scanning blocks")

	events, err := listener.re.GetEvents(ctx, start, end)
	if err != nil {
		return err
	}

	err = listener.processEvents(events)
	if err != nil {
		listener.logger.WithError(err).Error("failed to process events")
		//return err
	}

	err = listener.ds.ChainMeta.SaveLastHeight(ctx, listener.chainID, end)
	if err != nil {
		return err
	}

	return nil
}

func (listener *ExchangeListener) processEvents(events []*OrderEvent) error {
	listener.logger.Debugf("%+v\n", events)

	if listener.orderbook == nil {
		return nil
	}

	for _, event := range events {
		ctx := context.Background()

		var order *model.Order
		var orderHashErr error
		orderSignHashFound := false
		hashes := []string{event.Hash.String(), event.SellHash.String(), event.BuyHash.String()}

		for _, orderHash := range hashes {
			if orderHash == wyvern.NullAddress {
				continue
			}

			order, orderHashErr = listener.orderbook.GetBySignHash(ctx, orderHash)
			if orderHashErr == nil {
				orderSignHashFound = true
				break
			}

			listener.logger.
				WithField("order_hash", orderHash).
				WithError(orderHashErr).
				Error("failed to get order by hash")
		}

		if !orderSignHashFound {
			listener.logger.
				WithField("order_hashes", hashes).
				Error("failed to get order by hash")
			continue
		}

		switch event.Type {
		case OrderApproved:
			{
				listener.logger.
					WithField("hash", event.Hash.String()).
					WithField("event", "OrderApproved").
					Info("event received")
				return listener.orderbook.Approve(ctx, order)
			}
		case OrderCancelled:
			{
				listener.logger.
					WithField("hash", event.Hash.String()).
					WithField("event", "OrderCancelled").
					Info("event received")
				return listener.orderbook.Cancel(ctx, order)
			}
		case OrdersMatched:
			{
				logger := listener.logger.
					WithField("hash", event.Hash.String()).
					WithField("sell_hash", event.SellHash.String()).
					WithField("buy_hash", event.BuyHash.String()).
					WithField("maker", event.Maker.String()).
					WithField("taker", event.Taker.String()).
					WithField("maker", event.Maker.String()).
					WithField("event", "OrdersMatched")

				logger.Info("event received")

				var newOwner *model.Account

				if order.Side == wyvern.Sell {
					taker, err := listener.ds.Accounts.GetByAddress(ctx, event.Taker.String())
					if err != nil {
						logger.WithError(err).Error("failed to get new owner")
						return err
					}
					newOwner = taker
				}

				if order.Side == wyvern.Buy {
					maker, err := listener.ds.Accounts.GetByAddress(ctx, event.Maker.String())
					if err != nil {
						logger.WithError(err).Error("failed to get new owner")
						return err
					}
					newOwner = maker
				}

				return listener.orderbook.Process(ctx, order, newOwner)
			}
		}
	}

	return nil
}

func (listener *ExchangeListener) Start(errCh chan error) {
	listener.logger.Info("starting exchange listener")

	for range listener.t.C {
		listener.logger.Info("getting chain events")

		err := listener.waitEvents(context.Background())
		if err != nil {
			listener.logger.WithError(err).Error("failed to process events")
			continue
		}
	}
}

func (listener *ExchangeListener) Stop() error {
	listener.logger.Info("stopping exchange listener")
	listener.t.Stop()
	return nil
}

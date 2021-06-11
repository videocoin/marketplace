package listener

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

type ExchangeListener struct {
	logStep  uint64
	scanFrom uint64
	cli      *ethclient.Client
	re       *EventReader
}

func NewExchangeListener(logStep, scanFrom uint64, url string, contractAddress string) (*ExchangeListener, error) {
	cli, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	re, err := NewEventReader(cli, contractAddress)
	if err != nil {
		return nil, err
	}

	return &ExchangeListener{
		logStep:  logStep,
		scanFrom: scanFrom,
		cli:      cli,
		re:       re,
	}, nil
}

func (listener *ExchangeListener) Run(ctx context.Context, handler func([]*OrderEvent) error) error {
	// TODO: READ FROM DB
	// knownHeight should be a recovery point in case of service's crash
	knownHeight := uint64(0)
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
		return nil
	}

	end := start + listener.logStep
	if end > number {
		end = number
	}

	events, err := listener.re.GetEvents(ctx, start, end)
	if err != nil {
		return err
	}

	err = handler(events)
	if err != nil {
		return err
	}

	// TODO: WRITE TO DB
	knownHeight = end
	return nil
}

func (listener *ExchangeListener) headNumber(ctx context.Context) (uint64, error) {
	header, err := listener.cli.HeaderByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Uint64(), nil
}

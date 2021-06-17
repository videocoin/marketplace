package listener

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EventReader struct {
	ca  []common.Address
	cli *ethclient.Client
	pa  *Parser
}

func NewEventReader(cli *ethclient.Client, contractAddress string) (*EventReader, error) {
	ca := common.HexToAddress(contractAddress)

	pa := NewParser()
	if err := pa.RegisterABI(strings.NewReader(exchangeEventsABI)); err != nil {
		return nil, err
	}

	return &EventReader{
		ca:  []common.Address{ca},
		cli: cli,
		pa:  pa,
	}, nil
}

func (reader *EventReader) GetEvents(ctx context.Context, start, end uint64) ([]*OrderEvent, error) {
	logs, err := reader.cli.FilterLogs(ctx, ethereum.FilterQuery{
		Addresses: reader.ca,
		FromBlock: big.NewInt(int64(start)),
		ToBlock:   big.NewInt(int64(end)),
		Topics:    [][]common.Hash{{orderApprovedPartOne, orderApprovedPartTwo, orderCanceled, ordersMatched}},
	})
	if err != nil {
		return nil, err
	}

	events := make([]*OrderEvent, 0, len(logs))
	for i := range logs {
		var event *OrderEvent
		var eventType int
		switch logs[i].Topics[0] {
		case orderApprovedPartOne:
			eventType = OrderApproved
		case orderApprovedPartTwo:
			eventType = OrderApproved
		case orderCanceled:
			eventType = OrderCancelled
		case ordersMatched:
			eventType = OrdersMatched
		default:
			continue
		}
		event, err = reader.unpackEvent(ctx, eventType, logs[i].TxHash)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	return events, nil
}

func (reader *EventReader) unpackEvent(ctx context.Context, eventType int, txHash common.Hash) (*OrderEvent, error) {
	receipt, err := reader.cli.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		return nil, nil
	}
	switch eventType {
	case OrderApproved:
		one := reader.matchTypeEvent(orderApprovedPartOne, receipt.Logs)
		two := reader.matchTypeEvent(orderApprovedPartTwo, receipt.Logs)
		return reader.toOrderAprovedEvent(one, two)
	case OrderCancelled:
		log := reader.matchTypeEvent(orderCanceled, receipt.Logs)
		return reader.toOrderCanceledEvent(log)
	case OrdersMatched:
		log := reader.matchTypeEvent(ordersMatched, receipt.Logs)
		return reader.toOrdersMatchedEvent(log)
	}
	return nil, nil
}

func (reader *EventReader) toOrdersMatchedEvent(log *types.Log) (*OrderEvent, error) {
	event := ordersMatchedEvent{}
	if err := reader.pa.Unpack(&event, log); err != nil {
		return nil, err
	}
	return &OrderEvent{
		Type:  OrdersMatched,
		Hash:  event.SellHash,
		Maker: event.Maker,
		Taker: event.Taker,
	}, nil
}

func (reader *EventReader) toOrderCanceledEvent(log *types.Log) (*OrderEvent, error) {
	event := orderCanceledEvent{}
	if err := reader.pa.Unpack(&event, log); err != nil {
		return nil, err
	}
	return &OrderEvent{
		Type: OrderCancelled,
		Hash: event.Hash,
	}, nil
}

func (reader *EventReader) toOrderAprovedEvent(partOne, partTwo *types.Log) (*OrderEvent, error) {
	event1 := orderApprovedPartOneEvent{}
	event2 := orderApprovedPartTwoEvent{}
	if err := reader.pa.Unpack(&event1, partOne); err != nil {
		return nil, err
	}
	if err := reader.pa.Unpack(&event2, partTwo); err != nil {
		return nil, err
	}
	if event1.Hash != event2.Hash {
		return nil, fmt.Errorf("aproved order1 hash and aproved order2 hash doesn't match; buy=%s; sell=%s", event1.Hash.Hex(), event2.Hash.Hex())
	}
	return nil, nil
}

func (reader *EventReader) matchTypeEvent(topic common.Hash, logs []*types.Log) *types.Log {
	for _, log := range logs {
		if len(log.Topics) == 0 {
			// TODO: log
		}
		if log.Topics[0] == topic {
			return log
		}
	}
	return nil
}

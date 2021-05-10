package minter

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/videocoin/marketplace/internal/contracts/dev/nft"
)

type Minter struct {
	ca       common.Address
	cli      *ethclient.Client
	contract *nft.Nft1155
}

func NewMinter(url string, contractAddress string) (*Minter, error) {
	cli, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	ca := common.HexToAddress(contractAddress)
	contract, err := nft.NewNft1155(ca, cli)
	if err != nil {
		return nil, err
	}

	return &Minter{
		ca:       ca,
		cli:      cli,
		contract: contract,
	}, nil
}

func (m *Minter) ContractAddress() common.Address {
	return m.ca
}

func (m *Minter) Mint(ctx context.Context, to common.Address) error {
	return nil
}

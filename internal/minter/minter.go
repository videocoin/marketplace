package minter

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/videocoin/common/crypto"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/videocoin/marketplace/internal/contracts/dev/nft"
)

const (
	MintingGasLimit uint64 = 100000
)

type Minter struct {
	ca       common.Address
	cli      *ethclient.Client
	contract *nft.Nft1155
	opts     bind.TransactOpts
	mtx      sync.Mutex
}

func NewMinter(url string, contractAddress string, contractKeyFile string, contractKeyPass string) (*Minter, error) {
	cli, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	ca := common.HexToAddress(contractAddress)
	contract, err := nft.NewNft1155(ca, cli)
	if err != nil {
		return nil, err
	}

	key, err := crypto.DecryptKeyFile(contractKeyFile, contractKeyPass)
	if err != nil {
		return nil, err
	}

	return &Minter{
		ca:       ca,
		cli:      cli,
		contract: contract,
		opts:     *bind.NewKeyedTransactor(key.PrivateKey),
	}, nil
}

func (m *Minter) ContractAddress() common.Address {
	return m.ca
}

func (m *Minter) Mint(ctx context.Context, to common.Address, id *big.Int) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	opts := m.getCallOpts(ctx)
	balance, err := m.contract.BalanceOf(opts, to, id)
	if err != nil {
		return err
	}
	if balance.Cmp(big.NewInt(0)) != 0 {
		return fmt.Errorf("token with ID %s already exists", id.String())
	}

	txOpts := m.getTxOpts(ctx)
	tx, err := m.contract.Mint0(txOpts, to, id)
	if err != nil {
		return err
	}

	return m.waitMined(ctx, tx)
}

func (m *Minter) getCallOpts(ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{
		Context: ctx,
	}
}

func (m *Minter) getTxOpts(ctx context.Context) *bind.TransactOpts {
	return &bind.TransactOpts{
		GasLimit: MintingGasLimit,
		Context:  ctx,
		From:     m.opts.From,
		Signer:   m.opts.Signer,
	}
}

func (m *Minter) waitMined(ctx context.Context, tx *types.Transaction) error {
	receipt, err := bind.WaitMined(ctx, m.cli, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction %s failed", tx.Hash().String())
	}

	return nil
}

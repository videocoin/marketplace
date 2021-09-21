package minter

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/videocoin/marketplace/internal/contracts/dev/nft"
)

const (
	MintingGasLimit uint64 = 300000
	ZeroAddress     string = "0000000000000000000000000000000000000000"
)

type Minter struct {
	ca       common.Address
	cli      *ethclient.Client
	contract *nft.NFT721
	opts     bind.TransactOpts
	mtx      sync.Mutex
}

func NewMinter(url string, chainId uint64, contractAddress string, contractKey string, contractKeyPass string) (*Minter, error) {
	cli, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}

	ca := common.HexToAddress(contractAddress)
	contract, err := nft.NewNFT721(ca, cli)
	if err != nil {
		return nil, err
	}

	key, err := keystore.DecryptKey([]byte(contractKey), contractKeyPass)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt a key %s: %v", contractKey, err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, big.NewInt(int64(chainId)))
	if err != nil {
		return nil, fmt.Errorf("failed to create tx signer: %v", err)
	}

	return &Minter{
		ca:       ca,
		cli:      cli,
		contract: contract,
		opts:     *opts,
	}, nil
}

func (m *Minter) ContractAddress() common.Address {
	return m.ca
}

func (m *Minter) Mint(ctx context.Context, to common.Address, id *big.Int, uri string) (*types.Transaction, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// opts := m.getCallOpts(ctx)
	// _, err := m.contract.OwnerOf(opts, id)
	// if err == nil {
	// 	return nil, fmt.Errorf("token with ID %s already exists", id.String())
	// }
	// if err.Error() != "execution reverted: ERC721: owner query for nonexistent token" {
	// 	return nil, err
	// }

	txOpts := m.getTxOpts(ctx)
	tx, err := m.contract.Mint(txOpts, to, id, uri)
	if err != nil {
		return nil, err
	}

	return tx, m.waitMined(ctx, tx)
}

func (m *Minter) UpdateTokenURI(ctx context.Context, id *big.Int, uri string) (*types.Transaction, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// opts := m.getCallOpts(ctx)
	// _, err := m.contract.OwnerOf(opts, id)
	// if err == nil {
	// 	return nil, fmt.Errorf("token with ID %s already exists", id.String())
	// }
	// if err.Error() != "execution reverted: ERC721: owner query for nonexistent token" {
	// 	return nil, err
	// }

	txOpts := m.getTxOpts(ctx)
	tx, err := m.contract.UpdateTokenURI(txOpts, id, uri)
	if err != nil {
		return nil, err
	}

	return tx, m.waitMined(ctx, tx)
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

package wallet

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
)

type Wallet struct {
	address common.Address
	client  *ethclient.Client
}

func NewWallet(address common.Address, client *ethclient.Client) *Wallet {
	return &Wallet{address: address, client: client}
}

func (w *Wallet) Address() common.Address {
	return w.address
}
func (w *Wallet) GetBalance(ctx context.Context) (uint64, error) {
	lastBlock, err := w.client.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	balance, err := w.client.BalanceAt(ctx, w.address, lastBlock.Number())
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}
func (w *Wallet) GetTransactions(ctx context.Context) ([]types.Log, error) {
	lastBlock, err := w.client.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	logs, err := w.client.FilterLogs(ctx, ethereum.FilterQuery{
		BlockHash: nil,
		FromBlock: big.NewInt(0).Sub(lastBlock.Number(), big.NewInt(4000)),
		ToBlock:   nil,
		Addresses: []common.Address{w.Address()},
		Topics:    nil,
	})
	return logs, err
}

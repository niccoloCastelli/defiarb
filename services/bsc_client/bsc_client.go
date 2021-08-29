package bsc_client

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(rawUrl string, ctx context.Context) (*ethclient.Client, error) {
	if rawUrl == "" {
		rawUrl = "https://bsc-dataseed.binance.org/"
	}
	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	if chainId.Uint64() != 56 {
		return nil, errors.New("wrong chain id")
	}
	return client, nil
}

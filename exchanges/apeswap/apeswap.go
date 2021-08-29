package apeswap

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/exchanges"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake"
)

const (
	routerAddress = "0xC0788A3aD43d79aa53B09c2EaCc313A787d1d607"
	Name          = "apeswap"
)

func init() {
	exchanges.RegisterExchange(Name, NewRouter)
}

func NewRouter(client *ethclient.Client, address string) (exchanges.Router, error) {
	return pancake.NewRouterWithName(client, routerAddress, Name)
}

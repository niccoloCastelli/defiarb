package value

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/exchanges"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake"
)

const (
	routerAddress = "--"
	name          = "pandayield"
)

func init() {
	exchanges.RegisterExchange(name, NewRouter)
}

func NewRouter(client *ethclient.Client, address string) (exchanges.Router, error) {
	return pancake.NewRouter(client, routerAddress)
}

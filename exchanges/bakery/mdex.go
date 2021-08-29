package value

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/exchanges"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake"
)

const (
	routerAddress = "0xcde540d7eafe93ac5fe6233bee57e1270d3e330f"
	name          = "bakery"
)

func init() {
	exchanges.RegisterExchange(name, NewRouter)
}

func NewRouter(client *ethclient.Client, address string) (exchanges.Router, error) {
	return pancake.NewRouterWithName(client, routerAddress, name)
}

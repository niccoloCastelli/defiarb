package value

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/exchanges"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake"
)

const (
	routerAddress = "0x7DAe51BD3E3376B8c7c4900E9107f12Be3AF1bA8"
	name          = "mdex"
)

func init() {
	exchanges.RegisterExchange(name, NewRouter)
}

func NewRouter(client *ethclient.Client, address string) (exchanges.Router, error) {
	return pancake.NewRouterWithName(client, routerAddress, name)
}

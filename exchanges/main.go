package exchanges

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/storage/models"
	"github.com/pkg/errors"
	"math/big"
)

type RouterConstructor func(client *ethclient.Client, address string) (Router, error)

var registry = map[string]RouterConstructor{}

func RegisterExchange(name string, constructor RouterConstructor) {
	registry[name] = constructor
}

func New(name string, client *ethclient.Client, address string) (Router, error) {
	c, ok := registry[name]
	if !ok {
		return nil, errors.New("exchange not found")
	}
	return c(client, address)
}
func Registered() []string {
	ret := make([]string, 0, len(registry))
	for key := range registry {
		ret = append(ret, key)
	}
	return ret
}

type Router interface {
	Name() string
	GetRouterAddress() common.Address
	GetFactoryAddress() common.Address
	GetRouterContract() RouterContract
	GetAllPairs(ctx context.Context) ([]common.Address, error)
	GetAllPairsAsync(ctx context.Context) (int64, chan *common.Address, chan error)

	GetPairInfo(address common.Address) (*models.LiquidityPool, error)
	GetPairPrice(address common.Address, token0Amount int64, token1Amount int64) (*TradeInfo, error)
	GetReserves(address common.Address) (float64, float64, error)
}

type RouterContract interface {
	// GetAmountIn(opts *bind.CallOpts, amountOut *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error)
	// GetAmountOut(opts *bind.CallOpts, amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error)
	GetAmountsIn(opts *bind.CallOpts, amountOut *big.Int, path []common.Address) ([]*big.Int, error)
	GetAmountsOut(opts *bind.CallOpts, amountIn *big.Int, path []common.Address) ([]*big.Int, error)
	// Quote(opts *bind.CallOpts, amountA *big.Int, reserveA *big.Int, reserveB *big.Int) (*big.Int, error)
}

type TradeInfo struct {
	Token0Price float64
	Token1Price float64
}

func (t TradeInfo) String() string {
	return fmt.Sprintf("%f %f (slippage: %f)", t.Token0Price, t.Token1Price, t.Slippage())
}
func (t TradeInfo) Slippage() float64 {
	return (1/t.Token1Price - t.Token0Price) / t.Token0Price
}

package pancake

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/contracts/bep20"
	"github.com/niccoloCastelli/defiarb/exchanges"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake/contracts/factory"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake/contracts/pair"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake/contracts/router"
	"github.com/niccoloCastelli/defiarb/storage/models"
	"github.com/niccoloCastelli/defiarb/utils"
	"math"
	"math/big"
	"sync"
)

const (
	RouterAddress       = "0x05fF2B0DB69458A0750badebc4f9e13aDd608C7F"
	factoryAddress      = "0xBCfCcbde45cE874adCB698cC183deBcF17952812"
	Name                = "pancake"
	pairsChanBufferSize = 100
)

func init() {
	exchanges.RegisterExchange(Name, NewRouter)
}

func NewRouterWithName(client *ethclient.Client, address string, routerName string) (exchanges.Router, error) {
	if address == "" {
		address = RouterAddress
	}
	routerAddr := common.HexToAddress(address)
	routerContract, err := router.NewRouter(routerAddr, client)
	if err != nil {
		return nil, err
	}
	w, _ := routerContract.WETH(nil)
	fmt.Println("router: ", routerAddr.String(), w)
	factoryAddr, err := routerContract.Factory(nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("factory: ", factoryAddr.String())
	factoryContract, err := factory.NewFactory(factoryAddr, client)
	if err != nil {
		return nil, err
	}

	return &Router{
		name:        routerName,
		client:      client,
		routerAddr:  routerAddr,
		factoryAddr: factoryAddr,
		router:      routerContract,
		factory:     factoryContract,
	}, nil
}
func NewRouter(client *ethclient.Client, address string) (exchanges.Router, error) {
	if address == "" {
		address = RouterAddress
	}
	return NewRouterWithName(client, address, Name)
}

type Router struct {
	name        string
	client      *ethclient.Client
	routerAddr  common.Address
	factoryAddr common.Address
	router      *router.Router
	factory     *factory.Factory
}

func (r Router) Name() string {
	return r.name
}
func (r Router) GetRouterAddress() common.Address {
	return r.routerAddr
}
func (r Router) GetFactoryAddress() common.Address {
	return r.factoryAddr
}
func (r Router) GetRouterContract() exchanges.RouterContract {
	return r.router
}
func (r Router) GetAllPairsAsync(ctx context.Context) (int64, chan *common.Address, chan error) {
	pairAddressesChan := make(chan *common.Address, pairsChanBufferSize)
	errChan := make(chan error, 1)
	ctx, cancelFn := context.WithCancel(ctx)

	lengthB, err := r.factory.AllPairsLength(nil)
	if err != nil {
		errChan <- err
		pairAddressesChan <- nil
		cancelFn()
		return 0, pairAddressesChan, errChan
	}
	length := lengthB.Int64()

	go func() {
		defer func() {
			cancelFn()
			pairAddressesChan <- nil
		}()
		for i := int64(0); i < length; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				pairAddr, err := r.factory.AllPairs(nil, big.NewInt(i))
				if err != nil {
					errChan <- err
					return
				}
				pairAddressesChan <- &pairAddr
			}
		}
	}()

	return length, pairAddressesChan, errChan
}
func (r Router) GetAllPairs(ctx context.Context) ([]common.Address, error) {
	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	length, pairAddressesChan, errChan := r.GetAllPairsAsync(ctx)
	pairAddresses := make([]common.Address, 0, int(length))
	errs := []error{}

	wg := sync.WaitGroup{}
	wg.Add(cap(pairAddressesChan))
	go func() {
		i := 0
		for {
			select {
			case err := <-errChan:
				errs = append(errs, err)
				i++
				wg.Done()
			case pairAddr := <-pairAddressesChan:
				if pairAddr == nil {
					cancelFn()
					return
				}
				pairAddresses = append(pairAddresses, *pairAddr)
				i++
				wg.Done()
				if i%100 == 0 {
					fmt.Printf("%d/%d\n", i, length)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	<-ctx.Done()
	return pairAddresses, utils.WrapErrors(errs...)
}
func (r Router) GetPairInfo(address common.Address) (*models.LiquidityPool, error) {
	lp, err := pair.NewPair(address, r.client)
	if err != nil {
		return nil, err
	}
	name, err := lp.Name(nil)
	if err != nil {
		return nil, err
	}
	symbol, err := lp.Symbol(nil)
	if err != nil {
		return nil, err
	}
	token0Addr, err := lp.Token0(nil)
	if err != nil {
		return nil, err
	}
	token1Addr, err := lp.Token1(nil)
	if err != nil {
		return nil, err
	}
	reserves, err := lp.GetReserves(nil)
	if err != nil {
		return nil, err
	}
	_, err = lp.Decimals(nil)
	if err != nil {
		return nil, err
	}
	token0, err := bep20.NewBep20(token0Addr, r.client)
	if err != nil {
		return nil, err
	}
	token1, err := bep20.NewBep20(token1Addr, r.client)
	if err != nil {
		return nil, err
	}
	token0Symbol, err := token0.Symbol(nil)
	if err != nil {
		return nil, err
	}
	token1Symbol, err := token1.Symbol(nil)
	if err != nil {
		return nil, err
	}

	dbLp := models.LiquidityPool{
		Address:       address.String(),
		Name:          name,
		Symbol:        symbol,
		Description:   fmt.Sprintf("%s - %s", token0Symbol, token1Symbol),
		Exchange:      r.Name(),
		Token0Address: token0Addr.String(),
		Token1Address: token1Addr.String(),
		Token0Weight:  1,
		Token1Weight:  1,
		Reserve0:      utils.NewBigInt(reserves.Reserve0),
		Reserve1:      utils.NewBigInt(reserves.Reserve1),
	}

	return &dbLp, nil
}
func (r Router) GetPairPrice(address common.Address, token0Amount int64, token1Amount int64) (*exchanges.TradeInfo, error) {
	lp, err := pair.NewPair(address, r.client)
	if err != nil {
		return nil, err
	}

	token0Addr, err := lp.Token0(nil)
	if err != nil {
		return nil, err
	}
	token1Addr, err := lp.Token1(nil)
	if err != nil {
		return nil, err
	}
	token0, err := bep20.NewBep20(token0Addr, r.client)
	if err != nil {
		return nil, err
	}
	token1, err := bep20.NewBep20(token1Addr, r.client)
	if err != nil {
		return nil, err
	}
	token0Decimals, err := token0.Decimals(nil)
	if err != nil {
		return nil, err
	}
	token1Decimals, err := token1.Decimals(nil)
	if err != nil {
		return nil, err
	}

	token0Unit := new(big.Int).Mul(big.NewInt(token0Amount), big.NewInt(int64(math.Pow10(int(token0Decimals)))))
	token1Unit := new(big.Int).Mul(big.NewInt(token1Amount), big.NewInt(int64(math.Pow10(int(token1Decimals)))))

	/*
		reserves, err := lp.GetReserves(nil)
			if err != nil {
				return 0, err
			}
		// This is the tricky bit.
		// The reserve call returns the reserves for token0 and token1 in a sorted order.
		// This means we need to check if our token addresses are sorted or not and flip the reserves if they are not sorted.
		stoken0, _ := sortAddressess(token0Addr, token1Addr)
		if stoken0 != token0Addr {
			// We're not sorted, so the reserves need to be flipped to represent the actual reserves.
			reserves.Reserve0, reserves.Reserve1 = reserves.Reserve1, reserves.Reserve0
		}

		quote, err := r.router.Quote(nil, token0Unit, reserves.Reserve0, reserves.Reserve1)
		if err != nil {
			return 0, err
		}
		return utils.FormatAmount(quote, token1Decimals), nil*/

	amounts0Out, err := r.router.GetAmountsOut(nil, token0Unit, []common.Address{token0Addr, token1Addr})
	if err != nil {
		return nil, err
	}
	amounts1Out, err := r.router.GetAmountsOut(nil, token1Unit, []common.Address{token1Addr, token0Addr})
	if err != nil {
		return nil, err
	}
	new(big.Int).Quo(amounts0Out[1], token0Unit)
	return &exchanges.TradeInfo{
		Token0Price: utils.FormatAmount(new(big.Int).Quo(amounts0Out[1], big.NewInt(token0Amount)), token1Decimals),
		Token1Price: utils.FormatAmount(new(big.Int).Quo(amounts1Out[1], big.NewInt(token1Amount)), token0Decimals),
	}, nil
}
func (r Router) GetReserves(address common.Address) (float64, float64, error) {
	lp, err := pair.NewPair(address, r.client)
	if err != nil {
		return 0, 0, err
	}
	token0Addr, err := lp.Token0(nil)
	if err != nil {
		return 0, 0, err
	}
	token1Addr, err := lp.Token1(nil)
	if err != nil {
		return 0, 0, err
	}
	token0, err := bep20.NewBep20(token0Addr, r.client)
	if err != nil {
		return 0, 0, err
	}
	token1, err := bep20.NewBep20(token1Addr, r.client)
	if err != nil {
		return 0, 0, err
	}
	token0Decimals, err := token0.Decimals(nil)
	if err != nil {
		return 0, 0, err
	}
	token1Decimals, err := token1.Decimals(nil)
	if err != nil {
		return 0, 0, err
	}
	reserves, err := lp.GetReserves(nil)
	if err != nil {
		return 0, 0, err
	}
	return utils.FormatAmount(reserves.Reserve0, token0Decimals), utils.FormatAmount(reserves.Reserve1, token1Decimals), nil
}
func (r Router) Client() *ethclient.Client {
	return r.client
}

/*lastBlock, err := r.client.BlockByNumber(context.Background(), nil)
if err != nil {
	return 0, err
}
optsCurrentBlock := &bind.CallOpts{BlockNumber: lastBlock.Number()}
optsPreviousBlock := &bind.CallOpts{BlockNumber: new(big.Int).Sub(lastBlock.Number(), big.NewInt(1))}
cum0PreviousBlock, err := lp.Price0CumulativeLast(optsPreviousBlock)
if err != nil {
	return 0, err
}
cum1PreviousBlock, err := lp.Price1CumulativeLast(optsPreviousBlock)
if err != nil {
	return 0, err
}
cum0CurrentBlock, err := lp.Price0CumulativeLast(optsCurrentBlock)
if err != nil {
	return 0, err
}
cum1CurrentBlock, err := lp.Price1CumulativeLast(optsCurrentBlock)
if err != nil {
	return 0, err
}
price0 := new(big.Int).Set(new(big.Int).Sub(cum0CurrentBlock, cum0PreviousBlock))
price1 := new(big.Int).Set(new(big.Int).Sub(cum1CurrentBlock, cum1PreviousBlock))

fmt.Printf("%d (%d)\n", price0.Uint64(), price1.Uint64())*/

func sortAddressess(tkn0, tkn1 common.Address) (common.Address, common.Address) {
	token0Rep := big.NewInt(0).SetBytes(tkn0.Bytes())
	token1Rep := big.NewInt(0).SetBytes(tkn1.Bytes())

	if token0Rep.Cmp(token1Rep) > 0 {
		tkn0, tkn1 = tkn1, tkn0
	}

	return tkn0, tkn1
}

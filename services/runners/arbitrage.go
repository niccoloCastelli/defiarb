package runners

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/services/runners/contracts/arbitrage"
	"go.uber.org/atomic"
	"math/big"
	"os"
	"strings"
	"time"
)

const (
	sleepTime = time.Second * 15
)

type ArbitrageRoute struct {
	Router1 common.Address
	Router2 common.Address
	Token0  common.Address
	Token1  common.Address
	Token2  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	Amount2 *big.Int
	Key     string
}

type ArbitrageRunner struct {
	client          *ethclient.Client
	contractAddress common.Address
	isBusy          atomic.Bool
}

func NewArbitrageRunner(client *ethclient.Client, contractAddress common.Address) *ArbitrageRunner {
	return &ArbitrageRunner{client: client, contractAddress: contractAddress}
}

func (r *ArbitrageRunner) IsBusy() bool {
	return r.isBusy.Load()
}
func (r *ArbitrageRunner) callArbitrage(ctx context.Context, caller *arbitrage.Arbitrage, args ArbitrageRoute) error {
	ctx, cancelFn := context.WithTimeout(ctx, time.Second*5)
	defer cancelFn()
	privKey, err := crypto.LoadECDSA(os.Getenv("KEY_FILE"))
	if err != nil {
		return err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(56))
	if err != nil {
		return err
	}
	txn, err := caller.StartArbitrage(auth, args.Router1, args.Router2, args.Token0, args.Token1, args.Token2, args.Amount0, args.Amount1, args.Amount2)
	if err != nil {
		if strings.HasPrefix(err.Error(), "failed to estimate gas needed: gas required exceeds allowance") {
			auth.GasLimit = 4000000
			auth.GasPrice = big.NewInt(5000000000)
			txn, err = caller.StartArbitrage(auth, args.Router1, args.Router2, args.Token0, args.Token1, args.Token2, args.Amount0, args.Amount1, args.Amount2)
			if err != nil {
				fmt.Println("Contract call error: ", err.Error())
				return err
			}
		} else {
			fmt.Println("Contract call error: ", err.Error())
			return err
		}

	}
	tx, isPending, err := r.client.TransactionByHash(ctx, txn.Hash())
	if err != nil {
		return err
	}
	fmt.Println("contract called, txn hash: ", tx.Hash().String(), isPending, err)
	return nil
}
func (r *ArbitrageRunner) Run(ctx context.Context, inputChan chan ArbitrageRoute) error {
	caller, err := arbitrage.NewArbitrage(r.contractAddress, r.client)
	if err != nil {
		return err
	}
	go func() {
		var lastCall time.Time
		for {
			select {
			case args := <-inputChan:
				if r.IsBusy() {
					continue
				}
				if time.Now().Sub(lastCall) < time.Minute*5 {
					fmt.Println("Skip repetition", args.Key)
				}
				r.isBusy.Store(true)
				fmt.Println(args.Router1, args.Router2, args.Token1, args.Amount1.Uint64())
				if err := r.callArbitrage(ctx, caller, args); err != nil {
					fmt.Println("callArbitrage error: ", err.Error())
					lastCall = time.Now()
					time.Sleep(sleepTime)
					r.isBusy.Store(false)
					continue
				}
				time.Sleep(sleepTime)
				lastCall = time.Now()
				r.isBusy.Store(false)
				// router1 common.Address, router2 common.Address, token0 common.Address, token1 common.Address, token2 common.Address, amount0 *big.Int, amount1 *big.Int, amount2 *big.Int

			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

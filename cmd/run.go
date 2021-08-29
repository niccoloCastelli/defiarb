package cmd

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/contracts/bep20"
	"github.com/niccoloCastelli/defiarb/exchanges"
	_ "github.com/niccoloCastelli/defiarb/exchanges/apeswap"
	_ "github.com/niccoloCastelli/defiarb/exchanges/bakery"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake/contracts/pair"
	"github.com/niccoloCastelli/defiarb/services/bsc_client"
	"github.com/niccoloCastelli/defiarb/services/runners"
	"github.com/niccoloCastelli/defiarb/storage"
	"github.com/niccoloCastelli/defiarb/storage/models"
	"github.com/niccoloCastelli/defiarb/storage/price_graph"
	"github.com/niccoloCastelli/defiarb/utils"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// runCmd represents the scanTokens command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run arbitrage bot",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("scanTokens called")
		if conf == nil {
			return errors.New("config not found")
		}
		fmt.Println("scan called", walletAddr)
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()
		stopChan := make(chan os.Signal)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		client, err := bsc_client.NewClient(conf.NodeUrl, ctx)
		if err != nil {
			return errors.WithStack(err)
		}

		Db, err := storage.NewDb(&conf.Db)
		if err != nil {
			return errors.WithStack(err)
		}
		Db = Db.LogMode(false)
		if err := storage.Migrate(Db); err != nil {
			return errors.WithStack(err)
		}
		exchangesMap := map[string]exchanges.Router{}
		for _, exchangeName := range exchanges.Registered() {
			exchange, err := exchanges.New(exchangeName, client, "")
			if err != nil {
				return err
			}
			exchangesMap[exchangeName] = exchange
		}

		supportedTokens := models.Tokens{}
		if err := supportedTokens.GetAll(Db); err != nil {
			return err
		}
		supportedTokensAddrs := make([]string, len(supportedTokens))
		//startTokens := map[string]models.Token{}
		tokensMap := map[string]models.Token{}
		for i, tokenInfo := range supportedTokens {
			supportedTokensAddrs[i] = strings.ToLower(tokenInfo.Address)
			//token, err := bep20.NewBep20(common.HexToAddress(tokenInfo.Address), client)
			tokensMap[strings.ToLower(tokenInfo.Address)] = supportedTokens[i]
		}
		suppordedLps := models.LiquidityPools{}
		if err := suppordedLps.LoadSupportedPools(Db.LogMode(true)); err != nil {
			return err
		}
		priceChangesChan := make(chan models.LiquidityPool, len(suppordedLps))
		errChan := make(chan error)
		lock := sync.Mutex{}
		wg := sync.WaitGroup{}
		validLps := []models.LiquidityPool{}
		fmt.Println("Scanning supported LPs")
		scannedLps := atomic.Int32{}
		for _, lp := range suppordedLps {
			exchange, ok := exchangesMap[lp.Exchange]
			if !ok { //|| exchange.Name() == "value"
				continue
			}
			if err := lp.GetTokens(Db); err != nil {
				fmt.Println("Db error: ", err.Error())
				continue
			}
			wg.Add(1)
			go func(exchange exchanges.Router, lp models.LiquidityPool) {
				done := false
				defer func() {
					if !done {
						wg.Done()
						val := scannedLps.Inc()
						if val%100 == 0 {
							fmt.Printf("%d/%d\n", val, len(suppordedLps))
						}
					}
				}()
				swapsChan := make(chan *pair.PairSwap, len(suppordedLps))
				reserve0, reserve1, err := exchange.GetReserves(common.HexToAddress(lp.Address))
				if err != nil {
					fmt.Println("ERROR (GetReserves): ", lp.Exchange, lp.Description, err.Error())
					return
				}
				if lp.Token0.LoanAmount == 0 {
					token, err := bep20.NewBep20(common.HexToAddress(lp.Token0Address), client)
					if err != nil {
						return
					}
					decimals, err := token.Decimals(nil)
					if err != nil {
						return
					}
					lp.Token0.LoanAmount = int64(10 ^ int(decimals))
				}
				if lp.Token1.LoanAmount == 0 {
					token, err := bep20.NewBep20(common.HexToAddress(lp.Token1Address), client)
					if err != nil {
						return
					}
					decimals, err := token.Decimals(nil)
					if err != nil {
						return
					}
					lp.Token1.LoanAmount = int64(10 ^ int(decimals))
				}
				price, err := exchange.GetPairPrice(common.HexToAddress(lp.Address), lp.Token0.LoanAmount, lp.Token1.LoanAmount)
				if err != nil {
					fmt.Println("ERROR (GetPairPrice): ", lp.Exchange, lp.Description, err.Error())
					return
				}
				if reserve0 == 0 || reserve1 == 0 { //|| price.Slippage() > 0.4
					fmt.Println("Skip pair ", lp.Exchange, lp.Description, reserve0, reserve1)
					return
				}
				lpPair, err := pair.NewPair(common.HexToAddress(lp.Address), client)
				if err != nil {
					fmt.Println("ERROR (NewPair): ", lp.Exchange, lp.Description, err.Error())
					return
				}
				subscr, err := lpPair.WatchSwap(nil, swapsChan, nil, nil)
				if err != nil {
					fmt.Println("ERROR (WatchSwap): ", lp.Exchange, lp.Description, err.Error())
					return
				}
				lp.Token0Price = price.Token0Price
				lp.Token1Price = price.Token1Price
				fmt.Println("Subscribed swaps for ", lp.Name)

				lock.Lock()
				validLps = append(validLps, lp)
				lock.Unlock()
				wg.Done()
				done = true

				fmt.Printf("supported lp:[%s] [%s] %s reserves: %f %f\n", lp.Exchange, lp.Description, price.String(), reserve0, reserve1)
				defer subscr.Unsubscribe()
				for {
					select {
					case _ = <-swapsChan:
						priceChangesChan <- lp
					case err := <-subscr.Err():
						if err != nil {
							if err.Error() == "unexpected EOF" {
								continue
							}
							errChan <- err
						}
					case <-ctx.Done():
						return
					}
				}
			}(exchange, lp)

		}
		wg.Wait()

		writeClient, err := ethclient.Dial(conf.NodeUrl) // PROD
		if err != nil {
			return err
		}
		arbChan := make(chan runners.ArbitrageRoute, 0)
		arbRunner := runners.NewArbitrageRunner(writeClient, common.HexToAddress(conf.ContractAddress)) //PROD
		if err = arbRunner.Run(ctx, arbChan); err != nil {
			return err
		}
		writeClientChainId, err := writeClient.ChainID(ctx)
		if err != nil {
			return err
		}
		fmt.Println("ArbitrageRunner started", writeClientChainId.String())
		fmt.Printf("Watching %d pairs\n", len(validLps))

		priceGraph := price_graph.NewPriceGraph(validLps...)
		for _, token := range supportedTokens {
			priceGraph.ShortestPath(token.Address)
		}
		fmt.Println("Start listening...")
		for {
			select {
			case lp := <-priceChangesChan:
				exchange, ok := exchangesMap[lp.Exchange]
				if !ok {
					continue
				}
				go func(exchange exchanges.Router, lp models.LiquidityPool) {
					tradeInfo, err := exchange.GetPairPrice(common.HexToAddress(lp.Address), lp.Token0.LoanAmount, lp.Token1.LoanAmount)
					if err != nil {
						fmt.Println("ERROR (priceChangesChan)", err.Error())
						return
					}
					priceGraph.UpdatePrice(lp.Address, tradeInfo.Token0Price, tradeInfo.Token1Price)
					paths0, err := priceGraph.ShortestPath(lp.Token0Address)
					if err != nil {
						fmt.Println("ERROR (ShortestPath)", err.Error())
						return
					}
					paths1, err := priceGraph.ShortestPath(lp.Token1Address)
					if err != nil {
						fmt.Println("ERROR (ShortestPath)", err.Error())
						return
					}
					paths := append(paths0, paths1...)
					if len(paths) > 0 {
						for _, path := range paths {
							// fmt.Println("Opportunity")
							if (path.StartExchange == pancake.Name) && !arbRunner.IsBusy() && len(path.Tokens) > 1 {
								token1, ok := tokensMap[strings.ToLower(path.Tokens[1])]
								if !ok {
									continue
								}
								//m, _ := json.MarshalIndent(path, "", "  ")

								router1 := exchangesMap[path.Exchanges[1]].GetRouterAddress()
								router2 := common.Address{}
								token2 := common.Address{}
								if len(path.Exchanges) > 2 {
									router2 = exchangesMap[path.Exchanges[2]].GetRouterAddress()
									token2 = common.HexToAddress(path.Tokens[2])
								}

								amount1 := utils.GetAmount(token1.LoanAmount, int64(token1.Decimals))
								fmt.Println("Flash opportunity!", len(paths), path.Name, path.Weight, time.Now(), amount1.String(), path.Exchanges, path.Tokens) //, string(m)
								fmt.Println(router1, router2, token2)
								if conf.ExecuteTrades {
									arbChan <- runners.ArbitrageRoute{
										Key:     path.Key,
										Router1: router1,
										Router2: router2,
										Token0:  common.HexToAddress(path.Tokens[0]),
										Token1:  common.HexToAddress(path.Tokens[1]),
										Token2:  token2,
										Amount0: big.NewInt(int64(0)),
										Amount1: amount1,
										Amount2: big.NewInt(int64(0)),
									}
								}
							}
						}
					}
				}(exchange, lp)
				// fmt.Printf("price updated: [%s] [%s] %s\n", lp.Exchange, lp.Description, tradeInfo.String())
			case err := <-errChan:
				return err
			case <-stopChan:
				fmt.Println("---- STOP ----")
				return nil
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

/*
BUSD: https://bscscan.com/address/0xe9e7cea3dedca5984780bafc599bd69add087d56
USDT: https://bscscan.com/address/0x55d398326f99059ff775485246999027b3197955
BTC: https://bscscan.com/token/0x7130d2a12b9bcbfae4f2634d864a1ee1ce3ead9c
ETH: https://bscscan.com/token/0x2170ed0880ac9a755fd29b2688956bd959f933f8
DAI: https://bscscan.com/token/0x1af3f329e8be154074d8769d1ffa4ee058b1dbc3
DOT: https://bscscan.com/token/0x7083609fce4d1d8dc0c979aab8c869ea2c873402
XRP: https://bscscan.com/token/0x1d2f0da169ceb9fc7b3144628db156f3f6c60dbe
LINK: https://bscscan.com/token/0xf8a0bf9cf54bb92f17374d9e9a321e6a111a51bd
BAND: https://bscscan.com/token/0xad6caeb32cd2c308980a548bd0bc5aa4306c6c18
LTC: https://bscscan.com/token/0x4338665cbb7b2485a8855a139b75d5e34ab0db94
EOS: https://bscscan.com/token/0x56b6fb708fc5732dec1afc8d8556423a2edccbd6
BCH: https://bscscan.com/token/0x8ff795a6f4d97e7887c79bea79aba5cc76444adf
XTZ: https://bscscan.com/token/0x16939ef78684453bfdfb47825f8a5f714f12623a
ONT: https://bscscan.com/token/0xfd7b3a77848f1c2d67e05e54d78d174a0c850335
ADA: https://bscscan.com/token/0x3ee2200efb3400fabb9aacf31297cbdd1d435d47
ATOM: https://bscscan.com/token/0x0eb3a705fc54725037cc9e008bdede697f62f335
YFII: https://bscscan.com/token/0x7f70642d88cf1c4a3a7abb072b53b929b653eda5
ZEC: https://bscscan.com/token/0x1ba42e5193dfa8b03d15dd1b86a3113bbbef8eeb
*/

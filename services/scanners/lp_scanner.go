package scanners

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/jinzhu/gorm"
	"github.com/niccoloCastelli/defiarb/contracts/bep20"
	"github.com/niccoloCastelli/defiarb/contracts/pancake_pair"
	"github.com/niccoloCastelli/defiarb/storage/models"
	"github.com/pkg/errors"
	"math/big"
)

const (
	chunkSize          = 10000
	addrChanBufferSize = 20000
)

type LpScanner struct {
	client *ethclient.Client
	db     *gorm.DB
}

func NewLpScanner(client *ethclient.Client, db *gorm.DB) *LpScanner {
	return &LpScanner{client: client, db: db}
}

func (s *LpScanner) Scan(ctx context.Context, address common.Address) ([]*pancake_pair.Pair, error) {
	lastBlock, err := s.client.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ctx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()
	ret := []*pancake_pair.Pair{}
	lastBlockNumber := lastBlock.Number().Int64()
	addressesToScan := make(chan common.Address, addrChanBufferSize)
	logsChan := make(chan types.Log, addrChanBufferSize*5)
	errChan := make(chan error, 1)
	var startBlock int64 = 1

	// Get transaction logs for block chunks
	go func() {
		defer cancelFn()
		for i := lastBlockNumber; i > startBlock; i = i - chunkSize {
			select {
			case <-ctx.Done():
				return
			default:
				fromBlock := i - chunkSize
				if fromBlock < 0 {
					fromBlock = 0
				}
				fmt.Printf("Blocks: %d-%d \n", fromBlock, i)
				logs, err := s.client.FilterLogs(ctx, ethereum.FilterQuery{
					FromBlock: big.NewInt(fromBlock),
					ToBlock:   big.NewInt(i),
					Addresses: []common.Address{address}, //
				})
				if err != nil {
					errChan <- errors.WithStack(err)
					return
				}
				for _, log := range logs {
					logsChan <- log
				}
			}
		}
	}()

	// Get transaction receipts for logs
	go func() {
		for {
			select {
			case log := <-logsChan:
				tx, err := s.client.TransactionReceipt(ctx, log.TxHash)
				if err != nil {
					if jErr, ok := err.(rpc.Error); ok {
						switch jErr.ErrorCode() {
						case -32000:
							//fmt.Println("JsonRPC error", jErr.ErrorCode(), jErr.Error())
							continue
						default:
							fmt.Println("JsonRPC error", jErr.ErrorCode(), jErr.Error())
						}
					}
					errChan <- errors.WithStack(err)
					return
				}

				for _, txLog := range tx.Logs {
					addressesToScan <- txLog.Address
				}
			case <-ctx.Done():
				return
			}

		}
	}()

	savedTokens := map[string]common.Address{}
	savedLps := map[string]common.Address{}
	for {
		select {
		case err := <-errChan:
			return nil, err
		case addrToScan := <-addressesToScan:
			_, lpSaved := savedLps[addrToScan.String()]
			_, tokenSaved := savedTokens[addrToScan.String()]
			if lpSaved || tokenSaved {
				continue
			}
			_, saved, err := s.saveLp(addrToScan, address)
			if err != nil {
				return nil, err
			}
			if saved {
				savedLps[addrToScan.String()] = addrToScan
				continue
			}
			_, saved, err = s.saveToken(addrToScan)
			if err != nil {
				return nil, err
			}
			if saved {
				savedTokens[addrToScan.String()] = addrToScan
				continue
			}
			continue
		case <-ctx.Done():
			goto endloop
		}
	}
endloop:
	return ret, nil
}
func (s LpScanner) saveLp(addr common.Address, exchangeAddr common.Address) (*models.LiquidityPool, bool, error) {
	pair, err := pancake_pair.NewPair(addr, s.client)
	if err != nil {
		return nil, false, nil
	}
	name, err := pair.Name(nil)
	if err != nil {
		return nil, false, nil
	}
	symbol, err := pair.Symbol(nil)
	if err != nil {
		return nil, false, nil
	}
	token0Addr, err := pair.Token0(nil)
	if err != nil {
		return nil, false, nil
	}
	token1Addr, err := pair.Token1(nil)
	if err != nil {
		return nil, false, nil
	}

	dbLp := models.LiquidityPool{
		Address:       addr.String(),
		Name:          name,
		Symbol:        symbol,
		Exchange:      exchangeAddr.String(),
		Token0Address: token0Addr.String(),
		Token1Address: token1Addr.String(),
	}

	if factoryAddr, err := pair.Factory(nil); err == nil {
		dbLp.FactoryAddress = factoryAddr.String()
	}
	if err := s.db.Save(&dbLp).Error; err != nil {
		return nil, false, err
	}
	//fmt.Printf("AddLpTokenInfo(\"valueDefi\", \"%s\", \"%s\", \"%s\")\n", addr.String(), name, symbol)
	// fmt.Printf("[%s] %s %f %f\n", symbol, name, utils.FormatAmount(reserves.Reserve0), utils.FormatAmount(reserves.Reserve1))
	return &dbLp, true, nil
}
func (s LpScanner) saveToken(addr common.Address) (*models.Token, bool, error) {
	token, err := bep20.NewBep20(addr, s.client)
	if err != nil {
		return nil, false, nil
	}
	name, err := token.Name(nil)
	if err != nil {
		return nil, false, nil
	}
	symbol, err := token.Symbol(nil)
	if err != nil {
		return nil, false, nil
	}
	decimals, err := token.Decimals(nil)
	if err != nil {
		return nil, false, nil
	}
	totalSupply, err := token.TotalSupply(nil)
	if err != nil {
		return nil, false, err
	}

	dbToken := models.Token{
		Address:     addr.String(),
		Name:        name,
		Symbol:      symbol,
		LogoUrl:     "",
		Decimals:    int(decimals),
		TotalSupply: totalSupply.Int64(),
	}
	if err := s.db.Save(&dbToken).Error; err != nil {
		return nil, false, err
	}
	//fmt.Printf("AddLpTokenInfo(\"valueDefi\", \"%s\", \"%s\", \"%s\")\n", addr.String(), name, symbol)
	fmt.Printf("[%s] %s %f\n", symbol, name, dbToken.FormattedTotalSupply())
	return &dbToken, true, nil
}

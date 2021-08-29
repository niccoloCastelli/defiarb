package tokens

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/contracts/pancake_pair"
	"github.com/niccoloCastelli/defiarb/utils"
	"math/big"
	"time"
)

var lpTokens = map[string]LpToken{}

func init() {
	//AddLpTokenInfo("valueDefi", "0x8DD39f0a49160cDa5ef1E2a2fA7396EEc7DA8267", "vBSWAP/WBNB Value LP", "vBSWAPWBNB")
}

func NewLpTokenInfo(exchange string, address string, name string, code string, token0Address string, token1Address string) *LpToken {
	return &LpToken{
		Erc20:         &erc20{address: address, name: name, code: code},
		exchange:      exchange,
		token0Address: token0Address,
		token1Address: token1Address,
	}
}

type LpToken struct {
	Erc20
	exchange      string
	token0Address string
	token1Address string
}

func (t *LpToken) Tokens(client *ethclient.Client) ([]Erc20, error) {
	instance, err := pancake_pair.NewPair(common.HexToAddress(t.Address()), client)
	if err != nil {
		return nil, err
	}
	symbol, _ := instance.Symbol(nil)
	name, _ := instance.Name(nil)
	reserves, err := instance.GetReserves(nil)
	fmt.Println("reserves", reserves, symbol, name)
	return nil, err
}
func (t *LpToken) TransactionHistory(ctx context.Context, client *ethclient.Client) error {
	lastBlock, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		return err
	}
	lastBlockNumber := lastBlock.Number().Int64()
	pair, err := pancake_pair.NewPair(common.HexToAddress(t.Address()), client)
	if err != nil {
		return err
	}

	for i := lastBlockNumber; i > (lastBlockNumber - 20000); i = i - 5000 {
		fromBlock := uint64(i - 5000)
		toBlock := uint64(i)
		logs, err := pair.FilterSwap(&bind.FilterOpts{
			Start: fromBlock,
			End:   &toBlock,
		}, nil, nil)
		if err != nil {
			return err
		}

		for logs.Next() {
			if logs.Event == nil {
				break
			}
			block, err := client.BlockByNumber(ctx, big.NewInt(int64(logs.Event.Raw.BlockNumber)))
			if err != nil {
				return err
			}
			blockTime := time.Unix(int64(block.Time()), 0)
			if logs.Event.Amount0In.Uint64() != 0 {
				price := utils.FormatAmount(logs.Event.Amount0In, 18) / utils.FormatAmount(logs.Event.Amount1Out, 18)
				fmt.Println("swap 0 -> 1", price, utils.FormatAmount(logs.Event.Amount0In, 18), utils.FormatAmount(logs.Event.Amount1Out, 18), blockTime)
			} else {
				price := utils.FormatAmount(logs.Event.Amount0Out, 18) / utils.FormatAmount(logs.Event.Amount1In, 18)
				fmt.Println("swap 1 -> 0", price, utils.FormatAmount(logs.Event.Amount1In, 18), utils.FormatAmount(logs.Event.Amount0Out, 18), blockTime)
			}

		}
	}
	return nil
}
func AddLpTokenInfo(exchange string, address string, name string, code string) {
	lpTokens[address] = *NewLpTokenInfo(exchange, address, name, code, "", "")
}

func SupportedLpTokens() []LpToken {
	ret := make([]LpToken, 0, len(lpTokens))
	for _, tokenInfo := range lpTokens {
		ret = append(ret, tokenInfo)
	}
	return ret
}

package runners

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/services/runners/contracts/arbitrage"
	"github.com/niccoloCastelli/defiarb/utils"
	"go.uber.org/atomic"
	"math/big"
	"testing"
)

func newTestClient(ctx context.Context) (*ethclient.Client, error) {
	client, err := ethclient.Dial(utils.GetEnv("ETH_TEST_CLIENT", "http://127.0.0.1:8545"))
	//client, err := ethclient.Dial(utils.GetEnv("ETH_TEST_CLIENT", "https://bsc-dataseed.binance.org/"))
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("chain id: ", chainId)
	return client, nil
}

func TestCallArbitrage(t *testing.T) {
	ctx := context.Background()
	client, err := newTestClient(ctx)

	//contractAddress := common.HexToAddress("0x4a81cff73f1b8c6d94f50EDC08A4DEe7fbC109C6") //LOCALE
	contractAddress := common.HexToAddress("0x98Eb6aD42B9924A14450eaf6BEdC691ebd7E5791") //PROD
	if err != nil {
		t.Fatal(err)
	}
	caller, err := arbitrage.NewArbitrage(contractAddress, client)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		client          *ethclient.Client
		contractAddress common.Address
		isBusy          atomic.Bool
	}
	type args struct {
		ctx    context.Context
		caller *arbitrage.Arbitrage
		args   ArbitrageRoute
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				client:          client,
				contractAddress: contractAddress,
			},
			args: args{ //Flash opportunity! 6 DAI->WBNB->BUSD 1.002517542306942 2021-04-18 12:25:54.910438051 +0200 CEST m=+341.790852059 1000000000000000000 [pancake bakery pancake] [0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3 0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c 0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56]
				ctx:    ctx,
				caller: caller,
				args: ArbitrageRoute{
					Router1: common.HexToAddress("0xcde540d7eafe93ac5fe6233bee57e1270d3e330f"),
					Router2: common.HexToAddress("0x05fF2B0DB69458A0750badebc4f9e13aDd608C7F"),
					Token0:  common.HexToAddress("0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3"),
					Token1:  common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
					Token2:  common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"),
					Amount0: big.NewInt(0),
					Amount1: utils.GetAmount(1, 18),
					Amount2: big.NewInt(0),
				},
			},
		},
		/*{
			fields: fields{
				client:          client,
				contractAddress: contractAddress,
			},
			args:   args{
				ctx:    ctx,
				caller: caller,
				args:   ArbitrageRoute{
					Router1: common.HexToAddress("0x05fF2B0DB69458A0750badebc4f9e13aDd608C7F"),
					Router2: common.HexToAddress("0x05fF2B0DB69458A0750badebc4f9e13aDd608C7F"),
					Token0:  common.HexToAddress("0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3"),
					Token1:  common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
					Token2:  common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"),
					Amount0: big.NewInt(0),
					Amount1: utils.GetAmount(1,18),
					Amount2: big.NewInt(0),
				},
			},
		},
		{
			fields: fields{
				client:          client,
				contractAddress: contractAddress,
			},
			args:   args{
				ctx:    ctx,
				caller: caller,
				args:   ArbitrageRoute{
					Router1: common.HexToAddress("0xcde540d7eafe93ac5fe6233bee57e1270d3e330f"),
					Router2: common.HexToAddress("0x05fF2B0DB69458A0750badebc4f9e13aDd608C7F"),
					Token0:  common.HexToAddress("0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3"),
					Token1:  common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
					Token2:  common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"),
					Amount0: big.NewInt(0),
					Amount1: utils.GetAmount(1,18),
					Amount2: big.NewInt(0),
				},
			},
		},
		{
			fields: fields{
				client:          client,
				contractAddress: contractAddress,
			},
			args:   args{
				ctx:    ctx,
				caller: caller,
				args:   ArbitrageRoute{
					Router1: common.HexToAddress("0xcde540d7eafe93ac5fe6233bee57e1270d3e330f"),
					Router2: common.HexToAddress(""),
					Token0:  common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"),
					Token1:  common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
					Token2:  common.HexToAddress(""),
					Amount0: big.NewInt(0),
					Amount1: utils.GetAmount(1,18),
					Amount2: big.NewInt(0),
				},
			},
		},
		{
			fields: fields{
				client:          client,
				contractAddress: contractAddress,
			},
			args:   args{
				ctx:    ctx,
				caller: caller,
				args:   ArbitrageRoute{
					Router1: common.HexToAddress("0xb7e19a1188776f32e8c2b790d9ca578f2896da7c"),
					Router2: common.HexToAddress(""),
					Token0:  common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"),
					Token1:  common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
					Token2:  common.HexToAddress(""),
					Amount0: big.NewInt(0),
					Amount1: utils.GetAmount(1,18),
					Amount2: big.NewInt(0),
				},
			},
		},
		{ //Flash opportunity! 1 WBNB->ETH [pancake value] [0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c 0x2170Ed0880ac9A755fd29B2688956BD959F933F8]
			fields: fields{
				client:          client,
				contractAddress: contractAddress,
			},
			args:   args{
				ctx:    ctx,
				caller: caller,
				args:   ArbitrageRoute{
					Router1: common.HexToAddress("0xb7e19a1188776f32e8c2b790d9ca578f2896da7c"),
					Router2: common.HexToAddress(""),
					Token0:  common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
					Token1:  common.HexToAddress("0x2170Ed0880ac9A755fd29B2688956BD959F933F8"),
					Token2:  common.HexToAddress(""),
					Amount0: big.NewInt(0),
					Amount1: utils.GetAmount(1,18),
					Amount2: big.NewInt(0),
				},
			},
		},
		{ //Flash opportunity! 1 WBNB->ETH [pancake value] [0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c 0x2170Ed0880ac9A755fd29B2688956BD959F933F8]
			fields: fields{
				client:          client,
				contractAddress: contractAddress,
			},
			args:   args{
				ctx:    ctx,
				caller: caller,
				args:   ArbitrageRoute{
					Router1: common.HexToAddress("0xcde540d7eafe93ac5fe6233bee57e1270d3e330f"),
					Router2: common.HexToAddress(""),
					Token0:  common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"),
					Token1:  common.HexToAddress("0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56"),
					Token2:  common.HexToAddress(""),
					Amount0: big.NewInt(0),
					Amount1: utils.GetAmount(1,18),
					Amount2: big.NewInt(0),
				},
			},
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ArbitrageRunner{
				client:          tt.fields.client,
				contractAddress: tt.fields.contractAddress,
				isBusy:          tt.fields.isBusy,
			}
			if err := r.callArbitrage(tt.args.ctx, tt.args.caller, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("callArbitrage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

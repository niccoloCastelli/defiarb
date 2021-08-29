package pancake

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake/contracts/factory"
	"github.com/niccoloCastelli/defiarb/exchanges/pancake/contracts/router"
	"github.com/niccoloCastelli/defiarb/services/bsc_client"
	"github.com/niccoloCastelli/defiarb/storage/models"
	"reflect"
	"testing"
)

func newClient(t *testing.T) *ethclient.Client {
	ctx := context.Background()
	client, err := bsc_client.NewClient("ws://192.168.3.3/ws", ctx)
	if err != nil {
		t.Error(err)
	}
	return client
}

func TestNewRouter(t *testing.T) {
	type args struct {
		client  *ethclient.Client
		address string
	}
	type testData struct {
		name    string
		args    args
		want    *Router
		wantErr bool
	}

	client := newClient(t)

	tests := []testData{
		{
			args: args{
				client: client,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRouter(tt.args.client, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRouter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRouter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouter_GetAllPairs(t *testing.T) {
	type fields struct {
		client      *ethclient.Client
		routerAddr  common.Address
		factoryAddr common.Address
		router      *router.Router
		factory     *factory.Factory
	}
	client := newClient(t)

	tests := []struct {
		name    string
		fields  fields
		want    []common.Address
		wantErr bool
	}{
		{
			fields: fields{
				client:      client,
				routerAddr:  common.HexToAddress(RouterAddress),
				factoryAddr: common.HexToAddress(factoryAddress),
				router:      nil,
				factory:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewRouter(tt.fields.client, tt.fields.factoryAddr.String())
			if err != nil {
				t.Errorf("NewRouter error = %v", err)
			}
			got, err := r.GetAllPairs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllPairs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) < 10000 {
				t.Errorf("GetAllPairs() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouter_GetPairInfo(t *testing.T) {
	type fields struct {
		client      *ethclient.Client
		routerAddr  common.Address
		factoryAddr common.Address
		router      *router.Router
		factory     *factory.Factory
	}
	type args struct {
		address common.Address
	}
	client := newClient(t)
	router, err := NewRouter(client, RouterAddress)
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		fields  *Router
		args    args
		want    *models.LiquidityPool
		wantErr bool
	}{
		{
			fields: router.(*Router),
			args: args{
				address: common.HexToAddress("0xA527a61703D82139F8a06Bc30097cC9CAA2df5A6"),
			},
			want: &models.LiquidityPool{
				Address:       "0xA527a61703D82139F8a06Bc30097cC9CAA2df5A6",
				Name:          "Pancake LPs",
				Symbol:        "Cake-LP",
				Exchange:      "0x05fF2B0DB69458A0750badebc4f9e13aDd608C7F",
				Token0Address: "0x0E09FaBB73Bd3Ade0a17ECC321fD13a19e81cE82",
				Token1Address: "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",
				Token0Weight:  1,
				Token1Weight:  1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Router{
				client:      tt.fields.client,
				routerAddr:  tt.fields.routerAddr,
				factoryAddr: tt.fields.factoryAddr,
				router:      tt.fields.router,
				factory:     tt.fields.factory,
			}
			got, err := r.GetPairInfo(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPairInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPairInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouter_GetPairPrice(t *testing.T) {
	type fields struct {
		client      *ethclient.Client
		routerAddr  common.Address
		factoryAddr common.Address
		router      *router.Router
		factory     *factory.Factory
	}
	client := newClient(t)
	router, err := NewRouter(client, RouterAddress)
	if err != nil {
		t.Error(err)
	}
	type args struct {
		address common.Address
	}
	tests := []struct {
		name    string
		fields  *Router
		args    args
		want    float64
		wantErr bool
	}{
		{
			fields: router.(*Router),
			args: args{
				address: common.HexToAddress("0xA527a61703D82139F8a06Bc30097cC9CAA2df5A6"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Router{
				client:      tt.fields.client,
				routerAddr:  tt.fields.routerAddr,
				factoryAddr: tt.fields.factoryAddr,
				router:      tt.fields.router,
				factory:     tt.fields.factory,
			}
			got, err := r.GetPairPrice(tt.args.address, 0, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPairPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Errorf("GetPairPrice() got = %v", got)
			}
		})
	}
}

package tokens

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/niccoloCastelli/defiarb/contracts/bep20"
	"math"
)

type Erc20 interface {
	LogoUrl() string
	Code() string
	Name() string
	Address() string
	BalanceOf(client *ethclient.Client, account common.Address) (uint64, error)
	FormattedBalance(client *ethclient.Client, account common.Address) (float64, error)
}

func NewErc20(address string, name string, code string, logoUrl string) Erc20 {
	return &erc20{address: address, name: name, code: code, logoUrl: logoUrl}
}

type erc20 struct {
	address string
	name    string
	code    string
	logoUrl string
}

func (t erc20) LogoUrl() string {
	return t.logoUrl
}
func (t erc20) Code() string {
	return t.code
}
func (t erc20) Name() string {
	return t.name
}
func (t erc20) Address() string {
	return t.address
}
func (t erc20) BalanceOf(client *ethclient.Client, account common.Address) (uint64, error) {
	instance, err := bep20.NewBep20(common.HexToAddress(t.address), client)
	if err != nil {
		return 0, err
	}
	tokenBalance, err := instance.BalanceOf(nil, account)
	if err != nil {
		return 0, err
	}
	return tokenBalance.Uint64(), nil
}
func (t erc20) FormattedBalance(client *ethclient.Client, account common.Address) (float64, error) {
	instance, err := bep20.NewBep20(common.HexToAddress(t.address), client)
	balance, err := t.BalanceOf(client, account)
	if err != nil {
		return 0, err
	}
	decimals, err := instance.Decimals(nil)
	if err != nil {
		return 0, err
	}
	return float64(balance) / math.Pow(10, float64(decimals)), nil
}

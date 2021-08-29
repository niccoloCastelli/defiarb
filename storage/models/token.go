package models

import (
	"github.com/jinzhu/gorm"
	"github.com/niccoloCastelli/defiarb/tokens"
	"math"
)

type Token struct {
	BaseModel
	Address     string `gorm:"primary_key"`
	Name        string
	Symbol      string
	LogoUrl     string
	Decimals    int
	TotalSupply int64
	LoanAmount  int64 //`gorm:"-"`
}

func NewTokenFromERC20(t tokens.Erc20) *Token {
	return &Token{
		Address:     t.Address(),
		Name:        t.Name(),
		Symbol:      t.Code(),
		LogoUrl:     t.LogoUrl(),
		Decimals:    0,
		TotalSupply: 0,
		LoanAmount:  0,
	}
}
func (t Token) FormattedTotalSupply() float64 {
	return float64(t.TotalSupply) / math.Pow(10, float64(t.Decimals))
}

type Tokens []Token

func (t *Tokens) GetAll(db *gorm.DB) error {
	return db.Find(t).Error
}

package models

import (
	"github.com/jinzhu/gorm"
	"github.com/niccoloCastelli/defiarb/tokens"
	"github.com/niccoloCastelli/defiarb/utils"
)

type LiquidityPool struct {
	BaseModel
	Address        string `gorm:"primary_key"`
	FactoryAddress string
	Name           string
	Symbol         string
	Description    string
	Exchange       string
	Token0Address  string
	Token1Address  string
	Token0Weight   int
	Token1Weight   int
	Token0         Token
	Token1         Token
	Reserve0       *utils.BigInt `gorm:"type:bytea"`
	Reserve1       *utils.BigInt `gorm:"type:bytea"`
	Token0Price    float64       `gorm:"-"`
	Token1Price    float64       `gorm:"-"`
}

func (p *LiquidityPool) GetTokens(db *gorm.DB) error {
	if err := db.Where("lower(address) = lower(?)", p.Token0Address).First(&p.Token0).Error; err != nil {
		return err
	}
	if err := db.Where("lower(address) = lower(?)", p.Token1Address).First(&p.Token1).Error; err != nil {
		return err
	}
	return nil
}

const supportedPoolsQuery = `SELECT liquidity_pools.*, t0.symbol, t1.symbol FROM liquidity_pools
    INNER JOIN tokens t0 on liquidity_pools.token0_address = t0.address
    INNER JOIN tokens t1 on liquidity_pools.token1_address = t1.address`

type LiquidityPools []LiquidityPool

func (l *LiquidityPools) LoadSupportedPools(db *gorm.DB) error {
	return db.Raw(supportedPoolsQuery).Scan(&l).Error
}
func (l LiquidityPool) LpToken() *tokens.LpToken {
	return tokens.NewLpTokenInfo(l.Exchange, l.Address, l.Name, l.Symbol, l.Token0Address, l.Token1Address)
}

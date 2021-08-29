package utils

import (
	"math/big"
)

func FormatAmount(val *big.Int, decimals uint8) float64 {
	floatVal := new(big.Float).SetInt(val)
	exp := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	result, _ := new(big.Float).Quo(floatVal, new(big.Float).SetInt(exp)).Float64()
	return result
}

func GetAmount(baseAmount int64, decimals int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(baseAmount), new(big.Int).Exp(big.NewInt(10), big.NewInt(decimals), nil))
}

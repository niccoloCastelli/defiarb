package utils

import (
	"database/sql/driver"
	"fmt"
	"math/big"
)

type BigInt struct {
	big.Int
}

func NewBigInt(val *big.Int) *BigInt {
	if val == nil {
		return nil
	}
	return &BigInt{Int: *val}
}

func (b *BigInt) Value() (driver.Value, error) {
	if b != nil {
		return b.String(), nil
	}
	return nil, nil
}

func (b *BigInt) Scan(value interface{}) error {
	if value == nil {
		b = nil
	}
	switch t := value.(type) {
	case int64:
		b.Int = *big.NewInt(value.(int64))
	case []uint8:
		b.Int = *big.NewInt(0).SetBytes(t)
	default:
		return fmt.Errorf("Could not scan type %T into BigInt", t)
	}
	return nil
}

package main_test

import (
	"math/big"
	"time"

	"github.com/shopspring/decimal"
)

var Loc = time.FixedZone("Asia/Jakarta", int((7 * time.Hour).Seconds()))

func NewDecimal(nominal int64) decimal.Decimal {
	return decimal.NewFromBigInt(big.NewInt(nominal), -2)
}

func NewDate() time.Time {
	return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, Loc)
}

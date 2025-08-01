package models

import "github.com/shopspring/decimal"

type Currency struct {
	Coin      string          `json:"coin"`
	Timestamp uint32          `json:"timestamp"`
	Price     decimal.Decimal `json:"price"`
}

package dto

import "github.com/shopspring/decimal"

type PriceUSD struct {
	Usd decimal.Decimal `json:"usd"`
}

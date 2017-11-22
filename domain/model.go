package domain

import (
	"math/big"
	"time"
)

type Currency struct {
	Id   string
	Name string
}

type PortfolioEntry struct {
	CurrencyId string    `json:"currency"`
	Amount     big.Float `json:"amount"`
}

type Portfolio struct {
	Entries []*PortfolioEntry
}

type Quote struct {
	CurrencyId    string
	Price         big.Float
	Volume        big.Float
	High          big.Float
	Low           big.Float
	Change        string
	PercentChange string
	Period        time.Duration
}

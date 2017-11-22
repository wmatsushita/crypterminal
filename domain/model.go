package domain

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/wmatsushita/mycrypto/common"
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
	Change        big.Float
	PercentChange big.Float
	Period        time.Duration
}

func (q *Quote) String() string {
	quoteBytes, err := json.Marshal(q)
	if err != nil {
		fmt.Print(common.NewErrorWithCause("Could not marshal Quote to Json", err))
	}

	return string(quoteBytes)
}

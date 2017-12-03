package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/wmatsushita/crypterminal/common"
)

type Currency struct {
	Id   string
	Name string
}

type PortfolioEntry struct {
	CurrencyId string  `json:"currency"`
	Amount     float64 `json:"amount"`
}

type Portfolio struct {
	Entries []*PortfolioEntry
}

type Quote struct {
	CurrencyId    string
	Price         float64
	Volume        float64
	High          float64
	Low           float64
	Change        float64
	PercentChange float64
	Period        time.Duration
}

func (q *Quote) String() string {
	quoteBytes, err := json.Marshal(q)
	if err != nil {
		fmt.Print(common.NewErrorWithCause("Could not marshal Quote to Json", err))
	}

	return string(quoteBytes)
}

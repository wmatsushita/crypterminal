package bitfinex

import (
	"fmt"
	"net/http"
	"testing"
)

func TestBitfinexQuoteService_FetchQuotes(t *testing.T) {
	symbols := []string{"tBTCUSD", "tLTCUSD", "tIOTUSD"}

	service := NewBitfinexQuoteService(&http.Client{})

	quotes, err := service.FetchQuotes(symbols)
	if err != nil {
		t.Errorf("Fetch quotes returned error: %v", err)
	}

	for i, quote := range quotes {

		if quote.CurrencyId != symbols[i] {
			t.Errorf("Expected (%v), but found (%v)", symbols[i], quote.CurrencyId)
		}

		fmt.Printf("Quote: %v \n", quote)
	}
}

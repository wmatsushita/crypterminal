package bitfinex

import (
	"fmt"
	"net/http"
	"testing"
)

func TestBitfinexQuoteService_FetchQuotes(t *testing.T) {
	symbols := []string{"tBTCUSD", "tLTCUSD", "tIOTUSD"}

	service := NewBitfinexQuoteService(&http.Client{})

	fmt.Print(service.FetchQuotes(symbols))
}

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

	for currency, quote := range quotes {

		if quote.CurrencyId != currency {
			t.Errorf("Expected (%v), but found (%v)", currency, quote.CurrencyId)
		}

		fmt.Printf("Quote: %v \n", quote)
	}
}

func TestConvertionFromCurrencyToSymbol(t *testing.T) {
	currencyIds := []string{"Bitcoin", "Neo", "Iota", "Dash", "BitcoinGold"}
	symbols := convertFromCurrencyIdsToSymbols(currencyIds)

	expected := []string{"tBTCUSD", "tNEOUSD", "tIOTUSD", "tDSHUSD", "tBTGUSD"}

	for i, actual := range symbols {
		if actual != expected[i] {
			t.Errorf("Expected (%s) but got (%s)", expected[i], actual)
		}
	}
}

func TestAssembleQuoteRequest(t *testing.T) {
	Config = config{
		QuotesEndpoint: "http://mockendpoint.com",
		SymbolQueryKey: "dummyQuery",
	}

	request, err := assembleQuoteRequest([]string{"one", "two", "three"})
	if err != nil {
		t.Fatal(err)
	}

	actual, expected := request.URL.String(), "http://mockendpoint.com?dummyQuery=one%2Ctwo%2Cthree"
	if actual != expected {
		t.Errorf("Assembled request url (%s) differs from expected (%s).", actual, expected)
	}
}

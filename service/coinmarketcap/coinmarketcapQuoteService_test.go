package coinmarketcap

import (
	"fmt"
	"net/http"
	"testing"
)

func TestQuoteService_FetchQuotes(t *testing.T) {
	currencies := []string{"bitcoin", "litecoin", "iota"}

	service := NewCoinmarketcapQuoteService(&http.Client{}, "USD")

	quotes, err := service.FetchQuotes(currencies)
	if err != nil {
		t.Errorf("Fetch quotes returned error: %v", err)
	}

	if len(quotes) < 3 {
		t.Fatalf("Coinmarketcap response does not contain expected quotes")
	}

	for currency, quote := range quotes {

		if quote.CurrencyId != currency {
			t.Errorf("Expected (%v), but found (%v)", currency, quote.CurrencyId)
		}

		fmt.Printf("Quote: %v \n", quote)
	}
}

func TestAssembleQuoteRequest(t *testing.T) {

	QuoteEndpoint = "http://mockendpoint.com"

	service := NewCoinmarketcapQuoteService(&http.Client{}, "USD")

	request, err := service.assembleQuoteRequest()
	if err != nil {
		t.Fatal(err)
	}

	actual, expected := request.URL.String(), "http://mockendpoint.com?convert=USD&limit=200"
	if actual != expected {
		t.Errorf("Assembled request url (%s) differs from expected (%s).", actual, expected)
	}
}

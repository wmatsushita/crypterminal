package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"strconv"

	"github.com/wmatsushita/crypterminal/common"
	"github.com/wmatsushita/crypterminal/domain"
)

const (
	errorMsg string = "Could not fetch CoinMarketCap quotes"
)

type CoinmarketcapQuoteService struct {
	client       *http.Client
	fiatCurrency string
}

func NewCoinmarketcapQuoteService(c *http.Client, fiat string) *CoinmarketcapQuoteService {
	return &CoinmarketcapQuoteService{
		client:       c,
		fiatCurrency: fiat,
	}
}

func (s *CoinmarketcapQuoteService) FetchQuotes(currencyIds []string) (map[string]*domain.Quote, error) {

	request, err := s.assembleQuoteRequest()
	if err != nil {
		return errorResponse(errorMsg, err)
	}

	response, err := s.client.Do(request)
	if err != nil {
		return errorResponse(errorMsg, err)
	}

	bodyBytes, err := readResponseBody(response.Body)
	if err != nil {
		return errorResponse("Could not parse response", err)
	}

	if response.StatusCode != 200 {
		return nil, common.NewError(fmt.Sprintf("Received error response from Coinmarketcap: (%d)", response.StatusCode))
	}

	quotes, err := parseQuoteResponse(bodyBytes, currencyIds)
	if err != nil {
		return errorResponse(errorMsg, err)
	}

	return quotes, nil
}

func (s CoinmarketcapQuoteService) assembleQuoteRequest() (*http.Request, error) {
	url, err := url.Parse(QuoteEndpoint)
	if err != nil {
		return nil, common.NewErrorWithCause("Could not assemble Coinmarketcap request.", err)
	}

	query := url.Query()
	query.Set(FiatCurrencyParamKey, s.fiatCurrency)
	query.Set(LimitParamKey, LimitParamValue)
	url.RawQuery = query.Encode()

	return &http.Request{
		Method: http.MethodGet,
		URL:    url,
	}, nil
}

func parseQuoteResponse(body []byte, currencyIds []string) (map[string]*domain.Quote, error) {
	responseData := make([]QuoteResponse, 0)
	err := json.Unmarshal(body, &responseData)
	if err != nil {
		return nil, common.NewErrorWithCause("Could not Unmarshal server response", err)
	}

	quotes := make(map[string]*domain.Quote, 0)

	for _, record := range responseData {
		if common.Contains(currencyIds, record.Id) {
			quote := &domain.Quote{}

			quote.CurrencyId = record.Id
			quote.CurrencyName = record.Name
			quote.Price, err = strconv.ParseFloat(record.PriceUsd, 64)
			if err != nil {
				return nil, common.NewErrorWithCause("Could not convert quote price to number", err)
			}
			quote.Volume, err = strconv.ParseFloat(record.Volume24hUsd, 64)
			if err != nil {
				return nil, common.NewErrorWithCause("Could not convert quote volume to number", err)
			}
			quote.PercentChange, err = strconv.ParseFloat(record.PercentChange24h, 64)
			if err != nil {
				return nil, common.NewErrorWithCause("Could not convert quote percent change to number", err)
			}
			quote.PercentChange = quote.PercentChange / 100
			priceBefore := quote.Price / (quote.PercentChange + 1)
			quote.Change = quote.Price - priceBefore

			quote.Period = time.Hour * 24

			quotes[quote.CurrencyId] = quote
		}
	}

	return quotes, nil
}

func readResponseBody(body io.ReadCloser) ([]byte, error) {
	defer body.Close()

	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, common.NewErrorWithCause("Could not read response", err)
	}

	return bodyBytes, nil
}

func errorResponse(msg string, err error) (map[string]*domain.Quote, error) {
	return nil, common.NewErrorWithCause(msg, err)
}

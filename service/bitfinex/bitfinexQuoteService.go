package bitfinex

import (
	"net/http"

	"net/url"

	"strings"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"io"
	"time"

	"github.com/wmatsushita/mycrypto/common"
	"github.com/wmatsushita/mycrypto/domain"
)

const (
	errorMsg string = "Could not fetch Bitfinex quotes"
)

type BitfinexQuoteService struct {
	client *http.Client
}

func NewBitfinexQuoteService(c *http.Client) *BitfinexQuoteService {
	return &BitfinexQuoteService{
		client: c,
	}
}

func (s *BitfinexQuoteService) FetchQuotes(currencyIds []string) (map[string]*domain.Quote, error) {

	request, err := assembleQuoteRequest(convertFromCurrencyIdsToSymbols(currencyIds))
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
		return nil, common.NewError(fmt.Sprintf("Received error response from Bitfinex: (%d)", response.StatusCode))
	}

	quotes, err := parseQuoteResponse(bodyBytes)
	if err != nil {
		return errorResponse(errorMsg, err)
	}

	return quotes, nil
}

func convertFromCurrencyIdsToSymbols(currencyIds []string) []string {
	currencyToSymbol := getCurrencyToSymbolMap()
	symbols := make([]string, 0, len(currencyIds))

	for _, currencyId := range currencyIds {
		symbols = append(symbols, currencyToSymbol[currencyId])
	}

	return symbols
}

func assembleQuoteRequest(symbols []string) (*http.Request, error) {
	config := getConfig()
	fmt.Println(config)

	url, err := url.Parse(config.QuotesEndpoint)
	if err != nil {
		return nil, common.NewErrorWithCause("Could not assemble Bitfinex request.", err)
	}

	query := url.Query()
	query.Set(config.SymbolQueryKey, strings.Join(symbols, ","))
	url.RawQuery = query.Encode()

	return &http.Request{
		Method: http.MethodGet,
		URL:    url,
	}, nil
}

/*
Bitfinex response for ticker endpoint is as follows: (Nasty!)
[, col)
	[
		SYMBOL,				0
		BID,				1
		BID_SIZE,			2
		ASK,				3
		ASK_SIZE,			4
		DAILY_CHANGE,		5
		DAILY_CHANGE_PERC,	6
		LAST_PRICE,			7
		VOLUME,				8
		HIGH,				9
		LOW					10
	],
	...
]
*/
func parseQuoteResponse(body []byte) (map[string]*domain.Quote, error) {
	responseData := make([][]interface{}, 0)
	json.Unmarshal(body, &responseData)

	quotes := make(map[string]*domain.Quote, 0)

	for _, symbol := range responseData {
		quote := &domain.Quote{}
		quote.Period = time.Hour * 24

		for i, value := range symbol {
			switch i {
			case 0:
				quote.CurrencyId = value.(string)
			case 5:
				quote.Change = value.(float64)
			case 6:
				quote.PercentChange = value.(float64)
			case 7:
				quote.Price = value.(float64)
			case 8:
				quote.Volume = value.(float64)
			case 9:
				quote.High = value.(float64)
			case 10:
				quote.Low = value.(float64)
			}
		}
		quotes[quote.CurrencyId] = quote
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

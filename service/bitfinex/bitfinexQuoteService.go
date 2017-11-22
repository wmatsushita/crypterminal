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

	"math/big"

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

func (s *BitfinexQuoteService) FetchQuotes(symbols []string) ([]*domain.Quote, error) {

	request, err := assembleQuoteRequest(symbols)
	if err != nil {
		return errorResponse(errorMsg, err)
	}

	response, err := s.client.Do(request)
	if err != nil {
		return errorResponse(errorMsg, err)
	}

	quotes, err := parseQuoteResponse(response)
	if err != nil {
		return errorResponse(errorMsg, err)
	}

	return quotes, nil
}

func assembleQuoteRequest(symbols []string) (*http.Request, error) {
	config := getConfig()
	fmt.Println(config)

	url, err := url.Parse(config.QuotesEndpoint)
	if err != nil {
		return nil, common.NewErrorWithCause("Could not assemble Bitfinex request.", err)
	}

	url.Scheme = "https"
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
[
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
func parseQuoteResponse(response *http.Response) ([]*domain.Quote, error) {

	body, err := readResponseBody(response.Body)
	if err != nil {
		return errorResponse("Could not parse response", err)
	}

	if response.StatusCode != 200 {
		return nil, common.NewError(fmt.Sprintf("Received error response from Bitfinex: (%d)", response.StatusCode))
	}

	data := make([][]interface{}, 0)
	json.Unmarshal(body, &data)

	quotes := make([]*domain.Quote, 0)
	for _, row := range data {
		quote := &domain.Quote{}
		quote.Period = time.Hour * 24
		for i, col := range row {
			switch i {
			case 0:
				quote.CurrencyId = col.(string)
			case 5:
				convertToBigFloat(&quote.Change, col)
			case 6:
				convertToBigFloat(&quote.PercentChange, col)
			case 7:
				convertToBigFloat(&quote.Price, col)
			case 8:
				convertToBigFloat(&quote.Volume, col)
			case 9:
				convertToBigFloat(&quote.High, col)
			case 10:
				convertToBigFloat(&quote.Low, col)
			}
		}
		quotes = append(quotes, quote)
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

func convertToBigFloat(bigValue *big.Float, v interface{}) {
	value, ok := v.(float64)
	if ok {
		bigValue.SetFloat64(value)
	}
}

func errorResponse(msg string, err error) ([]*domain.Quote, error) {
	return nil, common.NewErrorWithCause(msg, err)
}

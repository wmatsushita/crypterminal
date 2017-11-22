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

	url, err := url.Parse(config.quotesEndpoint)
	if err != nil {
		return nil, common.NewErrorWithCause("Could not assemble Bitfinex request.", err)
	}

	url.Scheme = "https"
	query := url.Query()
	query.Set(config.symbolQueryKey, strings.Join(symbols, ","))
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

	data := make([]*quoteResponseRow, 0)
	json.Unmarshal(body, &data)

	quotes := make([]*domain.Quote, len(data))
	for _, q := range data {
		quote := &domain.Quote{}
		quote.Period = time.Hour * 24
		for i, col := range q.row {
			switch i {
			case 0:
				quote.CurrencyId = col
			case 5:
				quote.Change = col
			case 6:
				quote.PercentChange = col
			case 7:
				quote.Price.Parse(col, 10)
			case 8:
				quote.Volume.Parse(col, 10)
			case 9:
				quote.High.Parse(col, 10)
			case 10:
				quote.Low.Parse(col, 10)
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

type quoteResponseRow struct {
	row []string
}

func errorResponse(msg string, err error) ([]*domain.Quote, error) {
	return nil, common.NewErrorWithCause(msg, err)
}

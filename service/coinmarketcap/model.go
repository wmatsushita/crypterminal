package coinmarketcap

/*
Example data:
{
	"id": "bitcoin",
	"name": "Bitcoin",
	"symbol": "BTC",
	"rank": "1",
	"price_usd": "573.137",
	"price_btc": "1.0",
	"24h_volume_usd": "72855700.0",
	"market_cap_usd": "9080883500.0",
	"available_supply": "15844176.0",
	"total_supply": "15844176.0",
	"percent_change_1h": "0.04",
	"percent_change_24h": "-0.3",
	"percent_change_7d": "-0.57",
	"last_updated": "1472762067"
}
*/

type QuoteResponse struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUsd         string `json:"price_usd"`
	PriceBtc         string `json:"price_btc"`
	Volume24hUsd     string `json:"24h_volume_usd"`
	MarketCapUsd     string `json:"market_cap_usd"`
	AvailableSupply  string `json:"available_supply"`
	TotalSupply      string `json:"total_supply"`
	PercentChange1h  string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}

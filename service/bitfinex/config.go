package bitfinex

import (
	"log"

	"github.com/wmatsushita/crypterminal/common"
)

var (
	Config           config
	CurrencyToSymbol map[string]string
	SymbolToCurrency map[string]string
)

type config struct {
	QuotesEndpoint    string           `json:"quote.endpoint"`
	SymbolQueryKey    string           `json:"symbol.query.key"`
	CurrencySymbolMap []currencySymbol `json:"currency.symbol.map"`
}

type currencySymbol struct {
	CurrencyId string `json:"currency"`
	Symbol     string `json:"symbol"`
}

func getConfig() *config {
	if Config.QuotesEndpoint == "" {
		err := common.LoadFromJsonFile("bitfinex.config.json", &Config)
		if err != nil {
			log.Fatal(common.NewErrorWithCause("Could not load bitfinex.config.json file", err))
		}
	}

	return &Config
}

func getCurrencyToSymbolMap() map[string]string {
	cfg := getConfig()

	if len(CurrencyToSymbol) == 0 {
		CurrencyToSymbol = make(map[string]string)

		for _, cs := range cfg.CurrencySymbolMap {
			CurrencyToSymbol[cs.CurrencyId] = cs.Symbol
		}
	}

	return CurrencyToSymbol
}

func getSymbolToCurrencyMap() map[string]string {
	cfg := getConfig()

	if len(SymbolToCurrency) == 0 {
		SymbolToCurrency = make(map[string]string)

		for _, cs := range cfg.CurrencySymbolMap {
			SymbolToCurrency[cs.Symbol] = cs.CurrencyId
		}
	}

	return SymbolToCurrency
}

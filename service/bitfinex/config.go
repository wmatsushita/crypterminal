package bitfinex

import (
	"log"

	"github.com/wmatsushita/mycrypto/common"
)

var Config config

type config struct {
	QuotesEndpoint string `json:"quote.endpoint"`
	SymbolQueryKey string `json:"symbol.query.key"`
}

func getConfig() *config {
	if Config == (config{}) {
		err := common.LoadFromJsonFile("bitfinex.config.json", &Config)
		if err != nil {
			log.Fatal(common.NewErrorWithCause("Could not load bitfinex.config.json file", err))
		}
	}

	return &Config
}

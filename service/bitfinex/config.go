package bitfinex

import (
	"log"

	"github.com/wmatsushita/mycrypto/common"
)

var Config *config

type config struct {
	quotesEndpoint string `json:"quotesendpoint"`
	symbolQueryKey string `json:"symbolquerykey"`
}

func getConfig() *config {
	if Config == nil {
		Config = &config{}
		err := common.LoadFromJsonFile("bitfinex.config.json", Config)
		if err != nil {
			log.Fatal(common.NewErrorWithCause("Could not load bitfinex.config.json file", err))
		}
	}

	return Config
}

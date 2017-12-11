package main

import (
	"flag"

	"os"
	"os/signal"

	"net/http"

	"log"
	"time"

	"github.com/wmatsushita/crypterminal/mvp"
	"github.com/wmatsushita/crypterminal/service"
	"github.com/wmatsushita/crypterminal/service/bitfinex"
)

var (
	portfolioFlag    string
	serviceFlag      string
	fiatCurrencyFlag string
)

var httpClient *http.Client

func init() {
	flag.StringVar(&portfolioFlag, "p", "portfolio.json", "Portfolio filename, relative to current folder or absolute.")
	flag.StringVar(&serviceFlag, "service", "coinmarketcap", "API service to be used. Possible values are [coinmarketcap, bitfinex].")
	flag.StringVar(&fiatCurrencyFlag, "fiat", "USD", "Fiat currency used to show prices. Only works with coinmarketcap service. Possible values are [coinmarketcap, bitfinex].")

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
}

func main() {

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	portfolioService, err := service.NewJsonFilePortfolioService(portfolioFlag)
	if err != nil {
		log.Panicf("Error creating JsonFilePortfolioServcie: %s", err)
	}
	quoteService := bitfinex.NewBitfinexQuoteService(httpClient)

	quit := make(chan struct{}, 1)

	view := mvp.NewTermuiPortfolioView()
	presenter := mvp.NewPortfolioPresenter(view, quoteService, portfolioService, quit)

	presenter.Init()

	select {
	case <-interrupt:
		presenter.Quit()
		return
	case <-quit:
		return
	}

}

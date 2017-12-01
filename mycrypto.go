package main

import (
	"flag"

	"os"
	"os/signal"

	"net/http"

	"log"
	"time"

	"github.com/wmatsushita/mycrypto/mvp"
	"github.com/wmatsushita/mycrypto/service"
	"github.com/wmatsushita/mycrypto/service/bitfinex"
)

var (
	portfolioFlag string
	httpClient    *http.Client
)

func init() {
	flag.StringVar(&portfolioFlag, "portfolio", "portfolio.json", "Portfolio filename, relative to current folder or absolute.")

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
}

func main() {

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	view := mvp.NewTermuiPortfolioView()
	portfolioService, err := service.NewJsonFilePortfolioService(portfolioFlag)
	if err != nil {
		log.Panicf("Error creating JsonFilePortfolioServcie: %s", err)
	}
	quoteService := bitfinex.NewBitfinexQuoteService(httpClient)
	quit := make(chan struct{}, 1)
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

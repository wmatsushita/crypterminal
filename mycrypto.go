package main

import (
	"flag"

	"os"
	"os/signal"

	"github.com/wmatsushita/mycrypto/mvp"
)

var (
	portfolioFlag string
)

func init() {
	flag.StringVar(&portfolioFlag, "portfolio", "portfolio.json", "Portfolio filename, relative to current folder or absolute.")
}

func main() {

	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	view := mvp.NewTermuiPortfolioView()
	quit := make(chan struct{}, 1)
	presenter := mvp.NewPortfolioPresenter(view, quit)

	presenter.Init()

	select {
	case <-interrupt:
		presenter.Quit()
		return
	case <-quit:
		return
	}

}

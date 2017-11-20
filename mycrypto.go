package main

import (
	"flag"

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

	screen := mvp.NewTermuiPortfolioView()

	screen.Init()

}

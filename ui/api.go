package ui

import "github.com/wmatsushita/mycrypto/model"

type PortfolioScreen interface {
	Start()
	Quit()
	Refresh(portifolio []model.PortfolioEntry, quotes []model.Quote)
}

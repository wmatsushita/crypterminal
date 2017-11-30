package mvp

import (
	"fmt"

	"github.com/wmatsushita/mycrypto/domain"
)

const (
	FLOAT_FORMAT_STRING string = "%.4f"
)

type (
	Presenter interface {
		Init()
		ProcessUiEvent(event Event)
		Quit()
	}

	PortfolioPresenter struct {
		view PortfolioView
		quit chan struct{}
	}
)

func NewPortfolioPresenter(view PortfolioView, quit chan struct{}) *PortfolioPresenter {
	return &PortfolioPresenter{
		view,
		quit,
	}
}

func (p *PortfolioPresenter) Init() {
	p.view.Init(p)
}

func (p *PortfolioPresenter) Quit() {
	p.view.Quit()
	close(p.quit)
}

func (p *PortfolioPresenter) ProcessUiEvent(event Event) {
	switch event.Type {
	case programQuit:
		p.Quit()
	}
}

func (p *PortfolioPresenter) fillPortfolioTable(portfolio *domain.Portfolio, quotes map[string]*domain.Quote) {
	table := GetPortfolioTable()
	table.rows = make([]*PortfolioRow, 0, len(portfolio.Entries))
	for _, entry := range portfolio.Entries {
		quote := quotes[entry.CurrencyId]
		row := &PortfolioRow{
			assetName:     entry.CurrencyId,
			assetAmount:   formatValue(entry.Amount),
			assetPrice:    formatValue(quote.Price),
			assetValue:    formatValue(entry.Amount * quote.Price),
			valueChange:   formatValue(quote.Change),
			percentChange: formatValue(quote.PercentChange),
		}
		table.rows = append(table.rows, row)
	}

	// Notify observers that the table has been upated
	table.observable.Notify()
}

func formatValue(value float64) string {
	return fmt.Sprintf(FLOAT_FORMAT_STRING, value)
}

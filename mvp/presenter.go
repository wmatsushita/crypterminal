package mvp

import (
	"fmt"

	"time"

	"github.com/wmatsushita/crypterminal/domain"
	"github.com/wmatsushita/crypterminal/service"
)

const (
	tickInterval time.Duration = 10 * time.Second
	dateFormat   string        = "15:04:05"
)

var (
	ticker    <-chan time.Time
	portfolio *domain.Portfolio
)

type (
	Presenter interface {
		Init()
		ProcessUiEvent(event Event)
		Quit()
	}

	PortfolioPresenter struct {
		view             PortfolioView
		quoteService     service.QuoteService
		portfolioService service.PortfolioService
		quit             chan struct{}
		updated          bool
		lastUpdate       time.Time
	}
)

func NewPortfolioPresenter(
	view PortfolioView,
	quoteService service.QuoteService,
	portfolioService service.PortfolioService,
	quit chan struct{}) *PortfolioPresenter {

	return &PortfolioPresenter{
		view,
		quoteService,
		portfolioService,
		quit,
		false,
		time.Now(),
	}
}

func (p *PortfolioPresenter) Init() {
	p.view.Init(p)

	initializeTicker(p)

	p.reloadPortfolio()
}

func initializeTicker(p *PortfolioPresenter) {
	ticker = time.Tick(tickInterval)
	go func(t <-chan time.Time) {
		for range t {
			p.refreshQuotes()
		}
	}(ticker)
}

func (p *PortfolioPresenter) refreshQuotes() {
	p.setStatusMessage("Updating quotes...")
	p.fetchAndUpdateQuotes()
}

func (p *PortfolioPresenter) reloadPortfolio() {
	p.setStatusMessage("Reloading portfolio...")

	var err error
	portfolio, err = p.portfolioService.FetchPortfolio()
	if err != nil {
		p.setStatusMessage(fmt.Sprintf("Error reloading portfolio: %s", err))
	}

	p.fetchAndUpdateQuotes()
}

func (p *PortfolioPresenter) fetchAndUpdateQuotes() {
	quotes, err := p.fetchQuotesForPortfolio()
	if err != nil {
		msg := "Failed fetching quotes from server. "
		if p.updated {
			msg += p.lastUpdateMessage()
		}
		p.setStatusMessage(msg)
		return
	}

	p.updated = true
	p.lastUpdate = time.Now()
	p.fillPortfolioTable(portfolio, quotes)
	p.setStatusMessage(p.lastUpdateMessage())
}

func (p *PortfolioPresenter) lastUpdateMessage() string {
	return fmt.Sprintf("Last update: %v", p.lastUpdate.Format(dateFormat))
}

func (p *PortfolioPresenter) fetchQuotesForPortfolio() (map[string]*domain.Quote, error) {
	currencies := extractCurrencies(portfolio)
	quotes, err := p.quoteService.FetchQuotes(currencies)
	if err != nil {
		return nil, err
	}
	return quotes, nil
}

func extractCurrencies(portfolio *domain.Portfolio) []string {
	currencies := make([]string, 0, len(portfolio.Entries))
	for _, currency := range portfolio.Entries {
		currencies = append(currencies, currency.CurrencyId)
	}
	return currencies
}

func (p *PortfolioPresenter) setStatusMessage(msg string) {
	GetStatus().Msg = msg
	GetStatus().Observable.Notify()
}

func (p *PortfolioPresenter) Quit() {
	p.view.Quit()
	close(p.quit)
}

func (p *PortfolioPresenter) ProcessUiEvent(event Event) {
	switch event.Type {
	case portfolioRefresh:
		p.reloadPortfolio()
	case programQuit:
		p.Quit()
	}
}

func (p *PortfolioPresenter) fillPortfolioTable(portfolio *domain.Portfolio, quotes map[string]*domain.Quote) {
	table := GetPortfolioTable()
	table.Rows = make([]*PortfolioRow, 0, len(portfolio.Entries))
	for _, entry := range portfolio.Entries {
		quote := quotes[entry.CurrencyId]
		row := &PortfolioRow{
			AssetName:     quote.CurrencyName,
			AssetAmount:   entry.Amount,
			AssetPrice:    quote.Price,
			AssetValue:    entry.Amount * quote.Price,
			ValueChange:   quote.Change,
			PercentChange: quote.PercentChange * 100,
		}
		table.Rows = append(table.Rows, row)
	}

	// Notify observers that the table has been updated
	table.Observable.Notify()
}

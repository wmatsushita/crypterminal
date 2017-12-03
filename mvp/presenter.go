package mvp

import (
	"fmt"

	"time"

	"github.com/wmatsushita/mycrypto/domain"
	"github.com/wmatsushita/mycrypto/service"
)

const (
	FLOAT_FORMAT_STRING   string        = "%.4f"
	DOLAR_FORMAT_STRING   string        = "U$ %.2f"
	PERCENT_FORMAT_STRING string        = "%.2f %%"
	TICK_INTERVAL         time.Duration = 10 * time.Second
	DATE_FORMAT           string        = "15:04:05"
	ARROW_UP              string        = "\u21E7"
	ARROW_DOWN            string        = "\u21E9"
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
	}
}

func (p *PortfolioPresenter) Init() {
	p.view.Init(p)

	initializeTicker(p)

	p.reloadPortfolio()
}

func initializeTicker(p *PortfolioPresenter) {
	ticker = time.Tick(TICK_INTERVAL)
	go func(t <-chan time.Time) {
		for range t {
			p.refreshQuotes()
		}
	}(ticker)
}

func (p *PortfolioPresenter) refreshQuotes() {
	p.setStatusMessage("Updating quotes...")
	quotes, err := p.fetchQuotesForPortfolio()
	if err != nil {
		p.setStatusMessage("Failed to fetch quotes from server")
		return
	}

	p.fillPortfolioTable(portfolio, quotes)
	p.setStatusMessage(fmt.Sprintf("Last update: %v", time.Now().Format(DATE_FORMAT)))
}

func (p *PortfolioPresenter) reloadPortfolio() {
	p.setStatusMessage("Reloading portfolio...")

	var err error
	portfolio, err = p.portfolioService.FetchPortfolio()
	if err != nil {
		p.setStatusMessage(fmt.Sprintf("Error reloading portfolio: %s", err))
	}

	quotes, err := p.fetchQuotesForPortfolio()

	if err != nil {
		p.setStatusMessage("Failed fetching quotes from server")
		return
	}

	p.fillPortfolioTable(portfolio, quotes)
	p.setStatusMessage(fmt.Sprintf("Last update: %v", time.Now().Format(DATE_FORMAT)))
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
	totalValue := 0.0
	for _, entry := range portfolio.Entries {
		quote := quotes[entry.CurrencyId]
		row := &PortfolioRow{
			AssetName:     entry.CurrencyId,
			AssetAmount:   formatValue(FLOAT_FORMAT_STRING, entry.Amount),
			AssetPrice:    formatValue(DOLAR_FORMAT_STRING, quote.Price),
			AssetValue:    formatValue(DOLAR_FORMAT_STRING, entry.Amount*quote.Price),
			ValueChange:   formatChange(DOLAR_FORMAT_STRING, quote.Change),
			PercentChange: formatChange(PERCENT_FORMAT_STRING, quote.PercentChange*100),
		}
		table.Rows = append(table.Rows, row)
		totalValue += entry.Amount * quote.Price
	}

	totalRow := &PortfolioRow{
		AssetName:  "Total Portfolio Value:",
		AssetValue: formatValue(DOLAR_FORMAT_STRING, totalValue),
	}
	table.Rows = append(table.Rows, totalRow)

	// Notify observers that the table has been upated
	table.Observable.Notify()
}

func formatValue(format string, value float64) string {
	return fmt.Sprintf(format, value)
}

func formatChange(format string, change float64) string {
	if change > 0.0 {
		return ARROW_UP + " " + formatValue(format, change)
	} else {
		return ARROW_DOWN + " " + formatValue(format, change)
	}
}

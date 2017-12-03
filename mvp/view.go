package mvp

import (
	"time"

	"fmt"

	"github.com/gizak/termui"
	"github.com/wmatsushita/mycrypto/common"
)

const (
	amountFormatString  string = "%.4f"
	moneyFormatString   string = "$ %.2f"
	percentFormatString string = "%.2f %%"
	arrowUp             string = "\u21E7"
	arrowDown           string = "\u21E9"
)

type (
	PortfolioView interface {
		Init(presenter Presenter)
		Quit()
	}

	TermuiPortfolioView struct {
		presenter      Presenter
		title          *termui.Par
		menu           *termui.List
		portfolioTable *termui.Table
		statusBar      *termui.Par
		observer       common.Observer
		tableSignals   chan struct{}
		statusSignals  chan struct{}
		ticker         chan time.Time
	}
)

func NewTermuiPortfolioView() *TermuiPortfolioView {
	return &TermuiPortfolioView{
		title:          createTitle(),
		menu:           createMenu(),
		portfolioTable: createPortfolioTable(),
		statusBar:      createStatusBar(),
		observer:       common.NewEmptySignalObserver(),
	}
}

func (view *TermuiPortfolioView) Init(presenter Presenter) {
	view.presenter = presenter

	err := termui.Init()
	if err != nil {
		panic(err)
	}

	view.layout()
	view.eventHandling()

	view.tableSignals = view.observer.Watch(GetPortfolioTable().Observable, view.refreshPortfolioTable)
	view.statusSignals = view.observer.Watch(GetStatus().Observable, view.refreshStatus)
}

func (view *TermuiPortfolioView) refreshPortfolioTable() {
	data := GetPortfolioTable()
	view.portfolioTable.Rows = [][]string{{"Currency", "Ammount", "Price", "Value (USD)", "Value Change", "% Change"}}

	totalValue, totalValueChange := 0.0, 0.0
	for _, row := range data.Rows {
		view.portfolioTable.Rows = append(view.portfolioTable.Rows,
			[]string{
				row.AssetName,
				formatValue(amountFormatString, row.AssetAmount),
				formatValue(moneyFormatString, row.AssetPrice),
				formatValue(moneyFormatString, row.AssetValue),
				formatChange(moneyFormatString, row.ValueChange),
				formatChange(percentFormatString, row.PercentChange),
			})
		totalValue += row.AssetValue
		totalValueChange += row.AssetValue * row.PercentChange / 100
	}

	// Reset table row bgcolor to default
	for i := range view.portfolioTable.BgColors {
		view.portfolioTable.BgColors[i] = view.portfolioTable.BgColor
	}

	view.summaryRow(totalValueChange, totalValue)

	termui.Render(termui.Body)
}

func (view *TermuiPortfolioView) summaryRow(totalValueChange float64, totalValue float64) {
	totalPercentChange := totalValueChange / totalValue * 100
	view.portfolioTable.Rows = append(view.portfolioTable.Rows,
		[]string{
			"Total Portfolio Value",
			"", "",
			formatValue(moneyFormatString, totalValue),
			formatChange(moneyFormatString, totalValueChange),
			formatChange(percentFormatString, totalPercentChange),
		})
	// Set table last row (summary) bgcolor
	view.portfolioTable.BgColors[len(view.portfolioTable.Rows)-1] = termui.ColorBlue
}

func (view *TermuiPortfolioView) refreshStatus() {
	status := GetStatus()
	view.statusBar.Text = status.Msg
	termui.Render(termui.Body)
}

func (view *TermuiPortfolioView) eventHandling() {
	// handle key q pressing
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		view.presenter.ProcessUiEvent(Event{programQuit})
	})
	termui.Handle("/sys/kbd/r", func(termui.Event) {
		// press q to quit
		view.presenter.ProcessUiEvent(Event{portfolioRefresh})
	})

	go termui.Loop() // block until StopLoop is called
}

func (view *TermuiPortfolioView) layout() {
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(8, 0, view.title),
			termui.NewCol(4, 0, view.menu)),
		termui.NewRow(
			termui.NewCol(12, 0, view.portfolioTable),
		),
		termui.NewRow(
			termui.NewCol(12, 0, view.statusBar)))
	// calculate layout
	termui.Body.Align()
	termui.Render(termui.Body)
}

func (view *TermuiPortfolioView) Quit() {
	view.observer.Ignore(GetPortfolioTable().Observable, view.tableSignals)
	view.observer.Ignore(GetStatus().Observable, view.statusSignals)
	termui.StopLoop()
	termui.Close()
}

func createStatusBar() *termui.Par {
	statusBar := termui.NewPar("")
	statusBar.Height = 3
	statusBar.TextFgColor = termui.ColorGreen
	statusBar.BorderLabel = "Status"
	return statusBar
}

func createPortfolioTable() *termui.Table {
	tableData := [][]string{{"Currency", "Ammount", "Price", "Value (USD)", "Value Change", "% Change"}}

	table := termui.NewTable()
	table.Rows = tableData
	table.BorderLabel = "Portfolio"
	table.FgColor = termui.ColorWhite
	table.BgColor = termui.ColorDefault
	table.Height = 30
	table.Border = true
	table.FgColors = make([]termui.Attribute, 30)
	table.BgColors = make([]termui.Attribute, 30)
	table.FgColors[0] = termui.ColorCyan

	return table
}

func createMenu() *termui.List {
	strs := []string{
		"[q] Quit",
		"[r] Reload portfolio",
	}
	menu := termui.NewList()
	menu.Items = strs
	menu.ItemFgColor = termui.ColorYellow
	menu.BorderLabel = "Menu"
	menu.Height = 5

	return menu
}

func createTitle() *termui.Par {
	title := termui.NewPar(" \n   $$ MyCrypto Portfolio Ticker $$ ")
	title.Height = 5
	title.TextFgColor = termui.ColorGreen
	title.BorderFg = termui.ColorCyan

	return title
}

func formatValue(format string, value float64) string {
	return fmt.Sprintf(format, value)
}

func formatChange(format string, change float64) string {
	if change > 0.0 {
		return arrowUp + " " + formatValue(format, change)
	} else {
		return arrowDown + " " + formatValue(format, change)
	}
}

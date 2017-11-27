package mvp

import (
	"github.com/gizak/termui"
)

type (
	PortfolioView interface {
		Init(presenter Presenter)
		Watch(table *PortfolioTable)
		Quit()
	}

	TermuiPortfolioView struct {
		presenter      Presenter
		title          *termui.Par
		menu           *termui.List
		portfolioTable *termui.Table
		statusBar      *termui.Par
	}
)

func NewTermuiPortfolioView() *TermuiPortfolioView {
	return &TermuiPortfolioView{
		title:          createTitle(),
		menu:           createMenu(),
		portfolioTable: createPortfolioTable(),
		statusBar:      createStatusBar(),
	}
}

func (screen *TermuiPortfolioView) Init(presenter Presenter) {
	screen.presenter = presenter

	err := termui.Init()
	if err != nil {
		panic(err)
	}

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(10, 0, screen.title),
			termui.NewCol(2, 0, screen.menu)),
		termui.NewRow(
			termui.NewCol(12, 0, screen.portfolioTable),
		),
		termui.NewRow(
			termui.NewCol(12, 0, screen.statusBar)))

	// calculate layout
	termui.Body.Align()

	termui.Render(termui.Body)

	// handle key q pressing
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		// press q to quit
		presenter.ProcessUiEvent(Event{programQuit})
	})

	termui.Handle("/sys/kbd/r", func(termui.Event) {
		// press q to quit
		presenter.ProcessUiEvent(Event{portfolioRefresh})
	})

	go termui.Loop() // block until StopLoop is called
}

func (screen *TermuiPortfolioView) Watch(table *PortfolioTable) {
	panic("implement me")
}

func (screen *TermuiPortfolioView) Quit() {
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
	tableData := [][]string{
		{"Currency", "Ammount", "Price", "Daily Change"},
		{"Iota", "####", "$ ####.##", "100 %"}}
	table := termui.NewTable()
	table.Rows = tableData
	table.BorderLabel = "Portfolio"
	table.FgColor = termui.ColorWhite
	table.BgColor = termui.ColorDefault
	table.Height = 7
	table.Border = true

	return table
}
func createMenu() *termui.List {
	strs := []string{
		"[q] Quit",
		"[r] Reload portfolio input",
	}
	menu := termui.NewList()
	menu.Items = strs
	menu.ItemFgColor = termui.ColorYellow
	menu.BorderLabel = "Menu"
	menu.Height = 5

	return menu
}
func createTitle() *termui.Par {
	title := termui.NewPar(" \n   $$ MyCrypto Portifolio Ticker $$ ")
	title.Height = 5
	title.TextFgColor = termui.ColorGreen
	title.BorderFg = termui.ColorCyan

	return title
}

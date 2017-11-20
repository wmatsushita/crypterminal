package mvp

import (
	"github.com/gizak/termui"
)

type (
	PortfolioView interface {
		Init()
		Quit()
		Refresh(portifolio []PortfolioEntry, quotes []Quote)
	}

	TermuiPortfolioView struct {
		title          *termui.Par
		menu           *termui.List
		portfolioTable *termui.Table
		statusBar      *termui.Par
	}
)

func NewTermuiPortfolioScreen() *TermuiPortfolioView {
	return &TermuiPortfolioView{
		title:          createTitle(),
		menu:           createMenu(),
		portfolioTable: createPortfolioTable(),
		statusBar:      createStatusBar(),
	}
}

func (screen *TermuiPortfolioView) Init() {
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
		screen.Stop()
	})

	termui.Handle("/sys/kbd/C-x", func(termui.Event) {
		// handle Ctrl + x combination
	})

	termui.Handle("/sys/kbd", func(termui.Event) {
		// handle all other key pressing
	})

	// handle a 1s timer
	termui.Handle("/timer/1s", func(e termui.Event) {
		t := e.Data.(termui.EvtTimer)
		// t is a EvtTimer
		if t.Count%2 == 0 {
			// do something
		}
	})

	go termui.Loop() // block until StopLoop is called
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

func (screen *TermuiPortfolioView) Stop() {
	termui.StopLoop()
	termui.Close()
}

func (screen *TermuiPortfolioView) Refresh(portifolio []PortfolioEntry, quotes []Quote) {}

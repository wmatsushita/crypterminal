package mvp

import (
	"github.com/wmatsushita/mycrypto/common"
)

type (
	PortfolioRow struct {
		AssetName     string
		AssetAmount   string
		AssetPrice    string
		AssetValue    string
		ValueChange   string
		PercentChange string
	}

	PortfolioTable struct {
		Rows       []*PortfolioRow
		Observable common.Observable
	}

	Status struct {
		Msg        string
		Observable common.Observable
	}
)

var (
	thePortifolioTable = PortfolioTable{
		Rows:       make([]*PortfolioRow, 0),
		Observable: common.NewEmptySignalObservable(),
	}

	theStatus = Status{
		Observable: common.NewEmptySignalObservable(),
	}
)

func GetPortfolioTable() *PortfolioTable {
	return &thePortifolioTable
}

func GetStatus() *Status {
	return &theStatus
}

package mvp

import (
	"github.com/wmatsushita/crypterminal/common"
)

type (
	PortfolioRow struct {
		AssetName     string
		AssetAmount   float64
		AssetPrice    float64
		AssetValue    float64
		ValueChange   float64
		PercentChange float64
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

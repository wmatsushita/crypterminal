package mvp

import (
	"github.com/wmatsushita/mycrypto/common"
)

type (
	PortfolioRow struct {
		assetName     string
		assetAmount   string
		assetPrice    string
		assetValue    string
		valueChange   string
		percentChange string
	}

	PortfolioTable struct {
		rows       []*PortfolioRow
		observable common.Observable
	}

	Status struct {
		msg        string
		observable common.Observable
	}
)

var (
	thePortifolioTable = PortfolioTable{
		rows:       make([]*PortfolioRow, 0),
		observable: common.NewEmptySignalObservable(),
	}

	theStatus = Status{
		observable: common.NewEmptySignalObservable(),
	}
)

func GetPortfolioTable() *PortfolioTable {
	return &thePortifolioTable
}

func GetStatus() *Status {
	return &theStatus
}

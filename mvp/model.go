package mvp

import (
	"github.com/wmatsushita/mycrypto/common"
)

type (
	PortfolioRow struct {
		assetName    string
		assetAmount  string
		assetPrice   string
		assetValue   string
		valueChange  string
		percentChant string
	}

	PortfolioTable struct {
		rows       []PortfolioRow
		observable common.Observable
	}

	Status struct {
		msg        string
		observable common.Observable
	}
)

func NewPortfolioTable() *PortfolioTable {
	return &PortfolioTable{
		rows:       make([]PortfolioRow, 0),
		observable: common.NewEmptySignalObservable(),
	}
}

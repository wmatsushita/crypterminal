package service

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/wmatsushita/mycrypto/common"
	"github.com/wmatsushita/mycrypto/mvp"
)

type (
	PortfolioService interface {
		Load() error
		GetPortfolio() ([]*mvp.PortfolioEntry, error)
		common.Observable
	}

	JsonFilePortfolioService struct {
		filePath   string
		portfolio  []*mvp.PortfolioEntry
		observable *common.EmptySignalObservable
	}
)

func NewJsonFilePortfolioService(filePath string) (*JsonFilePortfolioService, error) {
	if !fileExists(filePath) {
		return nil, common.NewError("The given portifolio file path does not exist.")
	}
	return &JsonFilePortfolioService{
		filePath:   filePath,
		portfolio:  make([]*mvp.PortfolioEntry, 0),
		observable: common.NewEmptySignalObservable(),
	}, nil
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func (loader *JsonFilePortfolioService) Load() error {
	data, err := ioutil.ReadFile(loader.filePath)
	if err != nil {
		return common.NewErrorWithCause("Could not read portfolio file.", err)
	}

	err = json.Unmarshal(data, &loader.portfolio)
	if err != nil {
		return common.NewErrorWithCause("Could not deserialize json portfolio", err)
	}

	return nil
}

func (loader *JsonFilePortfolioService) GetPortfolio() ([]*mvp.PortfolioEntry, error) {
	if len(loader.portfolio) == 0 {
		return nil, common.NewError("The portfolio was not loaded. Consider loading it first.")
	}

	return loader.portfolio, nil
}

func (loader *JsonFilePortfolioService) Subscribe() chan struct{} {
	return loader.observable.Subscribe()
}

func (loader *JsonFilePortfolioService) Unsubscribe(subscription chan struct{}) {
	loader.observable.Unsubscribe(subscription)
}

func (loader *JsonFilePortfolioService) Notify() {
	loader.observable.Notify()
}

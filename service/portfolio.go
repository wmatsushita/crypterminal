package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/wmatsushita/mycrypto/common"
	"github.com/wmatsushita/mycrypto/mvp"
)

type (
	PortfolioLoader interface {
		Load() error
		GetPortfolio() ([]*mvp.PortfolioEntry, error)
		common.Observable
	}

	JsonFilePortfolioLoader struct {
		filePath   string
		portfolio  []*mvp.PortfolioEntry
		observable *common.EmptySignalObservable
	}
)

func NewJsonFilePortfolioLoader(filePath string) (*JsonFilePortfolioLoader, error) {
	if !fileExists(filePath) {
		return nil, common.NewError("The given portifolio file path does not exist.")
	}
	return &JsonFilePortfolioLoader{
		filePath:   filePath,
		portfolio:  make([]*mvp.PortfolioEntry, 0),
		observable: common.NewEmptySignalObservable(),
	}, nil
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func (loader *JsonFilePortfolioLoader) Load() error {
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

func (loader *JsonFilePortfolioLoader) GetPortfolio() ([]*mvp.PortfolioEntry, error) {
	if len(loader.portfolio) == 0 {
		return nil, common.NewError("The portfolio was not loaded. Consider loading it first.")
	}

	return loader.portfolio, nil
}

func (loader *JsonFilePortfolioLoader) Subscribe() chan struct{} {
	return loader.observable.Subscribe()
}

func (loader *JsonFilePortfolioLoader) Unsubscribe(subscription chan struct{}) {
	loader.observable.Unsubscribe(subscription)
}

func (loader *JsonFilePortfolioLoader) Notify() {
	loader.observable.Notify()
}

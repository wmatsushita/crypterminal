package main

import (
	"github.com/wmatsushita/mycrypto/mvp"
	"os"
	"encoding/json"
	"io/ioutil"
	"github.com/wmatsushita/mycrypto/common"
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
)

type (
	PortfolioLoader interface {
		Load() error
		GetPortfolio() ([]*mvp.PortfolioEntry, error)
		common.Observable
	}

	JsonFilePortfolioLoader struct {
		filePath  string
		portfolio []*mvp.PortfolioEntry
		subscriptions sets.Set
	}
)

func NewJsonFilePortfolioLoader(filePath string) (*JsonFilePortfolioLoader, error) {
	if !fileExists(filePath) {
		return nil, common.NewError("The given portifolio file path does not exist.")
	}
	return &JsonFilePortfolioLoader{
		filePath:  filePath,
		portfolio: make([]*mvp.PortfolioEntry,0),
		subscriptions: hashset.New(),
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

func (loader *JsonFilePortfolioLoader) Subscribe() chan bool {
	subscription := make(chan bool)
	loader.subscriptions.Add(subscription)

	return subscription
}

func (loader *JsonFilePortfolioLoader) Unsubscribe(subscription chan bool) {
	loader.subscriptions.Remove(subscription)
}

func (loader *JsonFilePortfolioLoader) Notify() {
	for  _, subscription := range loader.subscriptions.Values() {
		c, okToCast := subscription.(chan bool)
		if okToCast {
			go func() { c <- true }()
		}
	}
}

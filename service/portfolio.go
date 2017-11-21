package service

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/wmatsushita/mycrypto/common"
	"github.com/wmatsushita/mycrypto/domain"
)

type (
	PortfolioService interface {
		FetchPortfolio() (*domain.Portfolio, error)
	}

	JsonFilePortfolioService struct {
		filePath string
	}
)

func NewJsonFilePortfolioService(filePath string) (*JsonFilePortfolioService, error) {
	if !fileExists(filePath) {
		return nil, common.NewError("The given portifolio file path does not exist.")
	}
	return &JsonFilePortfolioService{
		filePath: filePath,
	}, nil
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func loadFromFile(filePath string) ([]*domain.PortfolioEntry, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, common.NewErrorWithCause("Could not read portfolio file.", err)
	}

	entries := make([]*domain.PortfolioEntry, 0)
	err = json.Unmarshal(data, &entries)
	if err != nil {
		return nil, common.NewErrorWithCause("Could not deserialize json portfolio", err)
	}

	return entries, nil
}

func (loader *JsonFilePortfolioService) FetchPortfolio() (*domain.Portfolio, error) {
	entries, err := loadFromFile(loader.filePath)
	if err != nil {
		return nil, err
	}

	return &domain.Portfolio{
		Entries: entries,
	}, nil
}

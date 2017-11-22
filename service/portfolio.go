package service

import (
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

func (loader *JsonFilePortfolioService) FetchPortfolio() (*domain.Portfolio, error) {
	entries := make([]*domain.PortfolioEntry, 0)
	err := common.LoadFromJsonFile(loader.filePath, &entries)
	if err != nil {
		return nil, err
	}

	return &domain.Portfolio{
		Entries: entries,
	}, nil
}

package service

import (
	"fmt"
	"testing"
)

func TestJsonPortfolioServiceValidatesNonExistingFiles(t *testing.T) {
	_, err := NewJsonFilePortfolioService("thisfiledoesnotexistforsure")
	if err == nil {
		t.Error("NewJsonFilePortfolioService did not fail when given unexisting file path.")
	}
}

func TestFileExistsValidation(t *testing.T) {
	fileName := "../portfolio.json"
	if !fileExists(fileName) {
		t.Errorf("fileExists function returns false for existing file name (%v).", fileName)
	}

	fileName = "thisfiledoesnotexistforsure986hlknasdoh"
	if fileExists(fileName) {
		t.Errorf("fileExsits function return true for non existing file name (%v)", fileName)
	}
}

func TestJsonPortfolioLoadsCorrectly(t *testing.T) {
	fileName := "../portfolio.json"

	loader, err := NewJsonFilePortfolioService(fileName)
	if err != nil {
		t.Error(err)
	}

	portfolio, err := loader.FetchPortfolio()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range portfolio.Entries {
		if v.CurrencyId == "IOTA" {
			iotaAmount := v.Amount.String()
			fmt.Printf("IOTA Amount: %v\n", iotaAmount)
			if iotaAmount != "20000" {
				t.Errorf("Iota amount difers from file content.")
			}
		}
	}

}

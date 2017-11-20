package main

import (
	"testing"
	"fmt"
)

func TestJsonPortfolioLoaderValidatesNonExistingFiles(t *testing.T) {
	_, err := NewJsonFilePortfolioLoader("thisfiledoesnotexistforsure")
	if err == nil {
		t.Error("NewJsonFilePortfolioLoader did not fail when given unexisting file path.")
	}
}

func TestFileExistsValidation(t *testing.T) {
	fileName := "portfolio.json"
	if !fileExists(fileName) {
		t.Errorf("fileExists function returns false for existing file name (%v).", fileName)
	}

	fileName = "thisfiledoesnotexistforsure986hlknasdoh"
	if fileExists(fileName) {
		t.Errorf("fileExsits function return true for non existing file name (%v)", fileName)
	}
}

func TestJsonPortfolioLoadsCorrectly(t *testing.T) {
	fileName := "portfolio.json"

	loader, err := NewJsonFilePortfolioLoader(fileName)
	if err != nil {
		t.Error(err)
	}

	err = loader.Load()
	if err != nil {
		t.Fatal(err)
	}

	portfolio, err := loader.GetPortfolio()
	if err != nil {
		t.Fatal(err)
	}

	for _,v := range portfolio {
		if v.CurrencyId == "IOTA" {
			iotaAmount := v.Amount.String()
			fmt.Printf("IOTA Amount: %v\n", iotaAmount)
			if iotaAmount != "20000" {
				t.Errorf("Iota amount difers from file content.")
			}
		}
	}

}
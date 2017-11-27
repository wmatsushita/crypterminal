package service

import "github.com/wmatsushita/mycrypto/domain"

type QuoteService interface {
	FetchQuotes(currencyIds []string) (map[string]*domain.Quote, error)
}

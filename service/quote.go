package service

import "github.com/wmatsushita/crypterminal/domain"

type QuoteService interface {
	FetchQuotes(currencyIds []string) (map[string]*domain.Quote, error)
}

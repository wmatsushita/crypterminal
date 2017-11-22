package service

import "github.com/wmatsushita/mycrypto/domain"

type QuoteService interface {
	FetchQuotes(symbols []string) ([]*domain.Quote, error)
}

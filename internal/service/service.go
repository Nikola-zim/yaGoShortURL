package service

import (
	"yaGoShortURL/internal/cash"
)

type CashURL interface {
	WriteURLInCash(string2 string) (string, error)
	ReadURLFromCash(string string) (string, error)
}

type Service struct {
	CashURL
}

func NewService(cash *cash.Cash) *Service {
	return &Service{
		CashURL: NewCashURLService(cash.UrlsRW),
	}
}

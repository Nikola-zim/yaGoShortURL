package service

import "yaGoShortURL/internal/app/cash"

type cashURL interface {
	WriteURLInCash(string2 string) (string, error)
	ReadURLFromCash(string string) (string, error)
}

type Service struct {
	cashURL
}

func NewService(cash *cash.Cash) *Service {
	return &Service{
		cashURL: NewCashURLService(cash.UrlsRW),
	}
}

package service

import (
	"yaGoShortURL/internal/cash"
)

type CashURLService struct {
	cash cash.UrlsRW
}

func NewCashURLService(cash cash.UrlsRW) *CashURLService {
	return &CashURLService{
		cash: cash,
	}
}

func (cu *CashURLService) WriteURLInCash(string2 string) (string, error) {
	return cu.cash.WriteURLInCash(string2)
}
func (cu *CashURLService) ReadURLFromCash(string string) (string, error) {
	return cu.cash.ReadURLFromCash(string)
}

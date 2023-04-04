package service

import (
	"yaGoShortURL/internal/cash"
)

type CashURLService struct {
	cash      cash.UrlsRW
	fileStore FileStoreURL
}

func NewCashURLService(cash cash.UrlsRW, fileStore FileStoreURL) *CashURLService {
	return &CashURLService{
		cash:      cash,
		fileStore: fileStore,
	}
}

func (cu *CashURLService) WriteURLInCash(fullURL string) (string, error) {
	id, err := cu.cash.WriteURLInCash(fullURL)
	if err != nil {
		return "", err
	}
	err = cu.fileStore.WriteURLInFile(fullURL, id)
	if err != nil {
		return "", err
	}
	return id, err
}
func (cu *CashURLService) ReadURLFromCash(string string) (string, error) {
	return cu.cash.ReadURLFromCash(string)
}

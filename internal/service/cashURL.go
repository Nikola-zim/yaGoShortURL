package service

import "encoding/binary"

type CashURLService struct {
	cash      CashURL
	fileStore FileStoreURL
}

func (cu *CashURLService) ReadAllUserURLFromCash(id []byte) ([]string, error) {
	return cu.cash.ReadAllUserURLFromCash(id)
}

func (cu *CashURLService) WriteURLInCash(fullURL string, userIDB []byte) (string, error) {
	id, err := cu.cash.WriteURLInCash(fullURL, userIDB)
	if err != nil {
		return "", err
	}
	userID := binary.LittleEndian.Uint64(userIDB)
	err = cu.fileStore.WriteURL(fullURL, id, userID)
	if err != nil {
		return "", err
	}
	return id, err
}

func (cu *CashURLService) ReadURLFromCash(string string) (string, error) {
	return cu.cash.ReadURLFromCash(string)
}

func NewCashURLService(cash CashURL, fileStore FileStoreURL) *CashURLService {
	return &CashURLService{
		cash:      cash,
		fileStore: fileStore,
	}
}

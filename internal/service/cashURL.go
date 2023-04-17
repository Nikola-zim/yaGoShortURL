package service

import "encoding/binary"

type CashURLService struct {
	cash      CashURL
	fileStore FileStoreURL
}

func (cu *CashURLService) ReadAllUserURLFromCash(id []byte) ([]string, error) {
	return cu.cash.ReadAllUserURLFromCash(id)
}

func (cu *CashURLService) WriteURLInCash(fullURL string, userIdB []byte) (string, error) {
	id, err := cu.cash.WriteURLInCash(fullURL, userIdB)
	if err != nil {
		return "", err
	}
	userId := binary.LittleEndian.Uint64(userIdB)
	err = cu.fileStore.WriteURL(fullURL, id, userId)
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

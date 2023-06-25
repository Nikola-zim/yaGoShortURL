package usecase

import (
	"yaGoShortURL/internal/entity"
)

type CashURLService struct {
	cash      CashURL
	fileStore FileStoreURL
	pg        DataBase
	usingDB   bool
}

func NewCashURLService(cash CashURL, fileStore FileStoreURL, pg DataBase, usingDB bool) *CashURLService {
	return &CashURLService{
		cash:      cash,
		fileStore: fileStore,
		pg:        pg,
		usingDB:   usingDB,
	}
}

func (cu *CashURLService) ReadAllUserURLFromCash(id uint64) ([]entity.JSONAllInfo, error) {
	return cu.cash.ReadAllUserURLFromCash(id)
}

func (cu *CashURLService) WriteURL(fullURL string, userID uint64) (string, error) {
	shortURL, err := cu.cash.WriteURL(fullURL, userID)
	if err != nil {
		return "", err
	}

	if cu.usingDB {
		err = cu.pg.WriteURL(fullURL, shortURL, userID)
	} else {
		err = cu.fileStore.WriteURLInFS(fullURL, shortURL, userID)
	}
	if err != nil {
		return "", err
	}
	return shortURL, err
}

func (cu *CashURLService) FullURL(string string) (string, error) {
	return cu.cash.FullURL(string)
}

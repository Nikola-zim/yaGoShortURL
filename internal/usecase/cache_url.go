package usecase

import (
	"yaGoShortURL/internal/entity"
)

type CashURLService struct {
	cash      CacheURL
	fileStore FileStoreURL
	pg        DataBase
	usingDB   bool
}

func NewCashURLService(cash CacheURL, fileStore FileStoreURL, pg DataBase, usingDB bool) *CashURLService {
	return &CashURLService{
		cash:      cash,
		fileStore: fileStore,
		pg:        pg,
		usingDB:   usingDB,
	}
}

func (cu *CashURLService) ReadAllUserURL(id uint64) ([]entity.JSONAllInfo, error) {
	return cu.cash.ReadAllUserURL(id)
}

func (cu *CashURLService) WriteURL(fullURL string, userID uint64) (string, error) {
	shortURL, err := cu.cash.WriteURL(fullURL, userID)
	if err != nil {
		return "", err
	}

	if cu.usingDB {
		err = cu.pg.WriteURL(fullURL, shortURL, userID)
	} else {
		err = cu.fileStore.WriteURL(fullURL, shortURL, userID)
	}
	if err != nil {
		return "", err
	}
	return shortURL, err
}

func (cu *CashURLService) FullURL(string string) (string, error) {
	return cu.cash.FullURL(string)
}

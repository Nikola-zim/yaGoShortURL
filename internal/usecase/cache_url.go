package usecase

import (
	"encoding/binary"
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

func (cu *CashURLService) ReadAllUserURLFromCash(id []byte) ([]entity.JSONAllInfo, error) {
	return cu.cash.ReadAllUserURLFromCash(id)
}

func (cu *CashURLService) WriteURL(fullURL string, userIDB []byte) (string, error) {
	id, err := cu.cash.WriteURL(fullURL, userIDB)
	if err != nil {
		return "", err
	}
	userID := binary.LittleEndian.Uint64(userIDB)

	if cu.usingDB {
		err = cu.pg.WriteURL(fullURL, id, userID)
	} else {
		err = cu.fileStore.WriteURLInFS(fullURL, id, userID)
	}
	if err != nil {
		return "", err
	}
	return id, err
}

func (cu *CashURLService) FullURL(string string) (string, error) {
	return cu.cash.FullURL(string)
}

package usecase

import (
	"context"
	"yaGoShortURL/internal/entity"
)

type CashURLService struct {
	cash        CacheURL
	fileStore   FileStoreURL
	pg          DataBase
	usingDB     bool
	toDeleteMsg chan entity.DeleteMsg
}

func NewCashURLService(cash CacheURL, fileStore FileStoreURL, pg DataBase, usingDB bool, msg chan entity.DeleteMsg) *CashURLService {
	return &CashURLService{
		cash:        cash,
		fileStore:   fileStore,
		pg:          pg,
		usingDB:     usingDB,
		toDeleteMsg: msg,
	}
}

func (cu *CashURLService) ReadAllUserURL(id uint64) ([]entity.JSONAllInfo, error) {
	return cu.cash.ReadAllUserURL(id)
}

func (cu *CashURLService) WriteURL(ctx context.Context, fullURL string, userID uint64) (string, error) {
	shortURL, err := cu.cash.WriteURL(fullURL, userID)
	if err != nil {
		return "", err
	}

	if cu.usingDB {
		err = cu.pg.WriteURL(ctx, fullURL, shortURL, userID)
	} else {
		err = cu.fileStore.WriteURL(fullURL, shortURL, userID)
	}
	if err != nil {
		return "", err
	}
	return shortURL, err
}

func (cu *CashURLService) FullURL(id string) (string, bool, error) {
	return cu.cash.FullURL(id)
}

func (cu *CashURLService) DeleteURLs(ctx context.Context, userID uint64, IDs []string) error {
	msg := entity.DeleteMsg{
		List:   IDs,
		UserID: userID,
	}
	cu.toDeleteMsg <- msg

	return nil
}

package usecase

import (
	"context"
	"log"
)

type DBService struct {
	cash CacheURL
	db   DataBase
}

func NewDBService(cash CacheURL, dataBase DataBase) *DBService {
	return &DBService{
		cash: cash,
		db:   dataBase,
	}
}

// PingDB Проверка связи с БД
func (m DBService) PingDB() error {
	return m.db.PingDB()
}

// RecoverAllURL Восстановление кэша из БД
func (m DBService) RecoverAllURL(ctx context.Context) error {
	dataURL, err := m.db.GetAllURL(ctx)
	if err != nil {
		log.Printf("internal - usecase - RecoverAllURL")
		return err
	}

	for _, data := range dataURL {
		_, err = m.cash.WriteURL(data.URL, data.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m DBService) WriteURLInDB(ctx context.Context, fullURL string, id string, userID uint64) error {
	err := m.db.WriteURL(ctx, fullURL, id, userID)
	if err != nil {
		log.Printf("internal - usecase - RecoverAllURL")
		return err
	}
	return nil
}

func (m DBService) DeleteURLsDB(ctx context.Context, userID uint64, IDs []string) error {
	return nil
}

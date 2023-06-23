package usecase

import (
	"context"
	"encoding/binary"
	"log"
)

type DBService struct {
	cash CashURL
	db   DataBase
}

func NewDBService(cash CashURL, dataBase DataBase) *DBService {
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
		userIDB := make([]byte, 8)
		binary.LittleEndian.PutUint64(userIDB, data.UserID)
		_, err = m.cash.WriteURL(data.URL, userIDB)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m DBService) WriteURLInDB(fullURL string, id string, userID uint64) error {
	err := m.db.WriteURL(fullURL, id, userID)
	if err != nil {
		log.Printf("internal - usecase - RecoverAllURL")
		return err
	}
	return nil
}

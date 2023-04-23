package service

import "yaGoShortURL/internal/static"

// CashURL интерфейс работы с кэшем
type CashURL interface {
	WriteURLInCash(fullURL string, userIDB []byte) (string, error)
	ReadURLFromCash(id string) (string, error)
	ReadAllUserURLFromCash(id []byte) ([]static.JSONAllInfo, error)
}

type AuthUser interface {
	FindUser(idMsg string) (uint64, bool)
	AddUser() (string, uint64, error)
}

// Cash Собранный интерфейс для кэша
type Cash interface {
	CashURL
	AuthUser
}

// FileStoreURL интерфейс работы с файлами
type FileStoreURL interface {
	WriteURL(fullURL string, id string, userID uint64) error
	ReadNextURL() (string, uint64, error)
}

type DB interface {
	PingDB() error
}

type Memory interface {
	RecoverAllURL() error
}

type Service struct {
	CashURLService
	MemoryService
	AuthService
	DBService
}

func NewService(cash Cash, fileStoreURL FileStoreURL, dataBase DB) *Service {
	return &Service{
		CashURLService: *NewCashURLService(cash, fileStoreURL),
		MemoryService:  *NewMemoryService(cash, fileStoreURL),
		AuthService:    *NewAuthService(cash),
		DBService:      *NewDBService(dataBase),
	}
}

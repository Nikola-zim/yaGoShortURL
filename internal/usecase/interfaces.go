package usecase

import (
	"context"
	"yaGoShortURL/internal/entity"
)

// CashURL интерфейс работы с кэшем
type CashURL interface {
	WriteURL(fullURL string, userIDB []byte) (string, error)
	FullURL(id string) (string, error)
	ReadAllUserURLFromCash(id []byte) ([]entity.JSONAllInfo, error)
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
	WriteURLInFS(fullURL string, id string, userID uint64) error
	ReadNextURL() (string, uint64, error)
}

type DataBase interface {
	PingDB() error
	WriteURL(fullURL string, id string, userID uint64) error
	GetAllURL(ctx context.Context) ([]entity.DataURL, error)
}

type Memory interface {
	RecoverAllURL() error
}

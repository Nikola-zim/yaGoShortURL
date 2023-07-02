package usecase

import (
	"context"
	"yaGoShortURL/internal/entity"
)

// CacheURL интерфейс работы с кэшем
type CacheURL interface {
	WriteURL(fullURL string, id uint64) (string, error)
	FullURL(id string) (string, bool, error)
	ReadAllUserURL(id uint64) ([]entity.JSONAllInfo, error)
	DeleteURLs(userID uint64, IDs []string) error
}

type AuthUser interface {
	FindUser(idMsg string) (uint64, bool)
	AddUser() (string, uint64, error)
}

// Cache Собранный интерфейс для кэша
type Cache interface {
	CacheURL
	AuthUser
}

// FileStoreURL интерфейс работы с файлами
type FileStoreURL interface {
	WriteURL(fullURL string, id string, userID uint64) error
	ReadNextURL() (string, uint64, error)
}

type DataBase interface {
	PingDB() error
	WriteURL(fullURL string, id string, userID uint64) error
	GetAllURL(ctx context.Context) ([]entity.DataURL, error)
	DeleteURLsDB(userID uint64, IDs []string) error
}

type Memory interface {
	RecoverAllURL() error
}

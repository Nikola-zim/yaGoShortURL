package workers

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

type DataBase interface {
	PingDB() error
	WriteURL(ctx context.Context, fullURL string, id string, userID uint64) error
	GetAllURL(ctx context.Context) ([]entity.DataURL, error)
	DeleteURLsDB(ctx context.Context, userID uint64, IDs []string) error
}

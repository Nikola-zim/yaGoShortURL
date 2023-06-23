package cache

import "yaGoShortURL/internal/entity"

type UrlsRW interface {
	WriteURL(fullURL string, id []byte) (string, error)
	FullURL(id string) (string, error)
	ReadAllUserURLFromCash(id []byte) ([]entity.JSONAllInfo, error)
}

type userGetAdd interface {
	FindUser(idMsg string) (uint64, bool)
	AddUser() (string, uint64, error)
}

type Cash struct {
	UrlsRW
	userGetAdd
}

func NewCash(baseURL string) *Cash {
	return &Cash{
		UrlsRW:     NewUrls(baseURL),
		userGetAdd: NewAuthUser(),
	}
}

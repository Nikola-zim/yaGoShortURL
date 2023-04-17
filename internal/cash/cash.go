package cash

import "yaGoShortURL/internal/static"

type UrlsRW interface {
	WriteURLInCash(fullURL string, id []byte) (string, error)
	ReadURLFromCash(id string) (string, error)
	ReadAllUserURLFromCash(id []byte) ([]static.JSONAllInfo, error)
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

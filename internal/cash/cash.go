package cash

import "yaGoShortURL/internal/static"

type UrlsRW interface {
	WriteURLInCash(fullURL string) (string, error)
	ReadURLFromCash(id string) (string, error)
}

type Cash struct {
	UrlsRW
	cfg static.ConfigInit
}

func NewCash(cfg static.ConfigInit) *Cash {
	return &Cash{
		UrlsRW: NewUrls(),
		cfg:    cfg,
	}
}

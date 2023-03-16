package cash

type UrlsRW interface {
	WriteURLInCash(fullURL string) (string, error)
	ReadURLFromCash(id string) (string, error)
}

type Cash struct {
	UrlsRW
}

func NewCash() *Cash {
	return &Cash{
		UrlsRW: NewUrls(),
	}
}

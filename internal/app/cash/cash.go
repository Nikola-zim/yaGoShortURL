package cash

type UrlsRW interface {
	WriteURLInCash(string2 string) (string, error)
	ReadURLFromCash(string string) (string, error)
}

type Cash struct {
	UrlsRW
}

func NewCash() *Cash {
	return &Cash{
		UrlsRW: NewUrls(),
	}
}

package cash

type UrlsRW interface {
	WriteUrlInCash(string2 string) error
	ReadUrlFromCash(int2 int) (string, error)
}

type Cash struct {
	UrlsRW
}

func NewCash() *Cash {
	return &Cash{
		UrlsRW: NewUrls(),
	}
}

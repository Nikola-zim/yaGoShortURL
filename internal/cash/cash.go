package cash

type UrlsRW interface {
	WriteURLInCash(fullURL string, id []byte) (string, error)
	ReadURLFromCash(id string) (string, error)
	ReadAllUserURLFromCash(id []byte) ([]string, error)
}

type userGetAdd interface {
	FindUser(idMsg string) (uint64, bool)
	AddUser() (string, error)
}

type Cash struct {
	UrlsRW
	userGetAdd
}

func NewCash() *Cash {
	return &Cash{
		UrlsRW:     NewUrls(),
		userGetAdd: NewAuthUser(),
	}
}

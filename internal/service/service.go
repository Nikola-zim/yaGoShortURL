package service

type CashURL interface {
	WriteURLInCash(fullURL string) (string, error)
	ReadURLFromCash(id string) (string, error)
}

// FileStoreURL Интерфейс работы с файлами
type FileStoreURL interface {
	WriteURL(fullURL string, id string) error
	ReadNextURL() (string, error)
}

type Memory interface {
	RecoverAllURL() error
}

type Service struct {
	CashURL
	Memory
}

func NewService(cash CashURL, fileStoreURL FileStoreURL) *Service {
	return &Service{
		CashURL: NewCashURLService(cash, fileStoreURL),
		Memory:  NewMemoryService(cash, fileStoreURL),
	}
}

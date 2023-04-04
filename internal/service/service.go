package service

type CashURL interface {
	WriteURLInCash(fullURL string) (string, error)
	ReadURLFromCash(id string) (string, error)
}

// FileStoreURL Интерфейс работы с файлами
type FileStoreURL interface {
	WriteURLInFile(fullURL string, id string) error
	ReadNextURLFromFile() (string, error)
}

type Memory interface {
	RecoverAllURL() error
}

type Service struct {
	CashURL
	Memory
	//FileStoreURL
}

func NewService(cash CashURL, fileStoreURL FileStoreURL) *Service {
	return &Service{
		CashURL: NewCashURLService(cash, fileStoreURL),
		Memory:  NewMemoryService(cash, fileStoreURL),
		//Передаю экземпляр кеша для записи всех URL в кеш при старте
		//FileStoreURL: NewFileStoreURLService(cash, fileStoreURL),
	}
}

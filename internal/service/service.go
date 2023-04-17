package service

// CashURL интерфейс работы с кэшем
type CashURL interface {
	WriteURLInCash(fullURL string, userIdB []byte) (string, error)
	ReadURLFromCash(id string) (string, error)
	ReadAllUserURLFromCash(id []byte) ([]string, error)
}

type AuthUser interface {
	FindUser(idMsg string) (uint64, bool)
	AddUser() (string, error)
}

// Cash Собранный интерфейс для кэша
type Cash interface {
	CashURL
	AuthUser
}

// FileStoreURL интерфейс работы с файлами
type FileStoreURL interface {
	WriteURL(fullURL string, id string, userID uint64) error
	ReadNextURL() (string, uint64, error)
}

type Memory interface {
	RecoverAllURL() error
}

type Service struct {
	CashURLService
	MemoryService
	AuthService
}

func NewService(cash Cash, fileStoreURL FileStoreURL) *Service {
	return &Service{
		CashURLService: *NewCashURLService(cash, fileStoreURL),
		MemoryService:  *NewMemoryService(cash, fileStoreURL),
		AuthService:    *NewAuthService(cash),
	}
}

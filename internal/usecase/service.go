package usecase

import "yaGoShortURL/internal/entity"

type Service struct {
	CashURLService
	MemoryService
	AuthService
	DBService
}

func NewService(cash Cache, fileStorage FileStoreURL, dataBase DataBase, usingDB bool, toDeleteMsg chan entity.DeleteMsg) *Service {
	return &Service{
		CashURLService: *NewCashURLService(cash, fileStorage, dataBase, usingDB, toDeleteMsg),
		MemoryService:  *NewMemoryService(cash, fileStorage),
		AuthService:    *NewAuthService(cash),
		DBService:      *NewDBService(cash, dataBase),
	}
}

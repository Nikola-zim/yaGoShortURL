package usecase

type Service struct {
	CashURLService
	MemoryService
	AuthService
	DBService
}

func NewService(cash Cache, fileStorage FileStoreURL, dataBase DataBase, usingDB bool) *Service {
	return &Service{
		CashURLService: *NewCashURLService(cash, fileStorage, dataBase, usingDB),
		MemoryService:  *NewMemoryService(cash, fileStorage),
		AuthService:    *NewAuthService(cash),
		DBService:      *NewDBService(cash, dataBase),
	}
}

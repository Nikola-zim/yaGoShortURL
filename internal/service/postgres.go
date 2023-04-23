package service

type DBService struct {
	dataBase DB
}

// PingDB Проверка связи с БД
func (m DBService) PingDB() error {
	return m.dataBase.PingDB()
}

func NewDBService(dataBase DB) *DBService {
	return &DBService{
		dataBase: dataBase,
	}
}

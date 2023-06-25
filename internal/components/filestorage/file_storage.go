package filestorage

type UrlsRW interface {
	WriteURL(fullURL string, id string, userID uint64) error
	ReadNextURL() (string, uint64, error)
	CloseFile() error
}

type FileStorage struct {
	UrlsRW
	unitTestFlag    bool
	fileStoragePath string
}

func NewFileStorage(unitTestFlag bool, fileStoragePath string) *FileStorage {
	//Отлавиливание ошибки при инициализации файлов
	UrlsRW, err := NewUrls(unitTestFlag, fileStoragePath)
	if err != nil {
		panic(err)
	}
	return &FileStorage{
		UrlsRW:          UrlsRW,
		unitTestFlag:    unitTestFlag,
		fileStoragePath: fileStoragePath,
	}
}

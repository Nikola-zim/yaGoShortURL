package filestorage

type UrlsRW interface {
	WriteURLInFile(fullURL string, id string) error
	ReadNextURLFromFile() (string, error)
	CloseFile() error
}

type FileStorage struct {
	UrlsRW
}

func NewFileStorage() *FileStorage {
	//Отлавиливание ошибки при инициализации файлов
	UrlsRW, err := NewUrls()
	if err != nil {
		panic(err)
	}
	return &FileStorage{
		UrlsRW: UrlsRW,
	}
}

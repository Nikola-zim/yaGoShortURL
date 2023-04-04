package fileStorage

type UrlsRW interface {
	WriteURLInFile(fullURL string, id string) error
	ReadAllURLFromFile(id string) (string, error)
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

package filestorage

import "yaGoShortURL/internal/static"

type UrlsRW interface {
	WriteURLInFile(fullURL string, id string) error
	ReadNextURLFromFile() (string, error)
	CloseFile() error
}

type FileStorage struct {
	UrlsRW
	cfg static.ConfigInit
}

func NewFileStorage(cfg static.ConfigInit) *FileStorage {
	//Отлавиливание ошибки при инициализации файлов
	UrlsRW, err := NewUrls(cfg)
	if err != nil {
		panic(err)
	}
	return &FileStorage{
		UrlsRW: UrlsRW,
		cfg:    cfg,
	}
}

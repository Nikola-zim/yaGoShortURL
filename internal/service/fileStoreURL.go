package service

import (
	"log"
	"yaGoShortURL/internal/cash"
)

type FileStoreURLService struct {
	cash      cash.UrlsRW
	fileStore FileStoreURL
}

func (f FileStoreURLService) WriteURLInFile(fullURL string, id string) error {
	log.Println(fullURL, id)
	return nil
}

func (f FileStoreURLService) ReadNextURLFromFile() (string, error) {
	return "", nil
}

//func NewFileStoreURLService(cash cash.UrlsRW, fileStore FileStoreURL) *FileStoreURLService {
//	return &FileStoreURLService{
//		cash:      cash,
//		fileStore: fileStore,
//	}
//}

package service

import "yaGoShortURL/internal/cash"

type FileStoreURLService struct {
	cash      cash.UrlsRW
	fileStore FileStoreURL
}

func (f FileStoreURLService) WriteURLInFile(fullURL string, id string) error {
	panic("implement me")
}

func (f FileStoreURLService) ReadAllURLFromFile(string string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewFileStoreURLService(cash cash.UrlsRW, fileStore FileStoreURL) *FileStoreURLService {
	return &FileStoreURLService{
		cash:      cash,
		fileStore: fileStore,
	}
}

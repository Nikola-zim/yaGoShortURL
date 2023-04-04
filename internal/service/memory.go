package service

import "yaGoShortURL/internal/cash"

type memoryService struct {
	cash      cash.UrlsRW
	fileStore FileStoreURL
}

func (m memoryService) RecoverAllURL() ([]string, error) {
	res := make([]string, 5)
	res[0], _ = m.fileStore.ReadAllURLFromFile("1")
	return res, nil
}

func NewMemoryService(cash cash.UrlsRW, fileStore FileStoreURL) *memoryService {
	return &memoryService{
		cash:      cash,
		fileStore: fileStore,
	}
}

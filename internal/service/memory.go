package service

import (
	"yaGoShortURL/internal/cash"
)

type MemoryService struct {
	cash      CashURL
	fileStore FileStoreURL
}

func (m MemoryService) RecoverAllURL() error {
	for {
		var nextURL string
		nextURL, err := m.fileStore.ReadNextURL()
		if err != nil || nextURL == "" {
			break
		}
		_, err = m.cash.WriteURLInCash(nextURL)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewMemoryService(cash cash.UrlsRW, fileStore FileStoreURL) *MemoryService {
	return &MemoryService{
		cash:      cash,
		fileStore: fileStore,
	}
}

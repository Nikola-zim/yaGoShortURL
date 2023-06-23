package usecase

import "encoding/binary"

type MemoryService struct {
	cash      CashURL
	fileStore FileStoreURL
}

func NewMemoryService(cash CashURL, fileStore FileStoreURL) *MemoryService {
	return &MemoryService{
		cash:      cash,
		fileStore: fileStore,
	}
}

// RecoverAllURL Восстановление кэша из файла
func (m MemoryService) RecoverAllURL() error {
	for {
		var url string
		url, userID, err := m.fileStore.ReadNextURL()
		if err != nil || url == "" {
			break
		}
		userIDB := make([]byte, 8)
		binary.LittleEndian.PutUint64(userIDB, userID)
		_, err = m.cash.WriteURL(url, userIDB)
		if err != nil {
			return err
		}
	}
	return nil
}

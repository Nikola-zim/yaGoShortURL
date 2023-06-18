package usecase

import "encoding/binary"

type MemoryService struct {
	cash      CashURL
	fileStore FileStoreURL
}

// RecoverAllURL Восстановление кэша из файла
func (m MemoryService) RecoverAllURL() error {
	for {
		var nextURL string
		nextURL, userID, err := m.fileStore.ReadNextURL()
		if err != nil || nextURL == "" {
			break
		}
		userIDB := make([]byte, 8)
		binary.LittleEndian.PutUint64(userIDB, userID)
		_, err = m.cash.WriteURL(nextURL, userIDB)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewMemoryService(cash CashURL, fileStore FileStoreURL) *MemoryService {
	return &MemoryService{
		cash:      cash,
		fileStore: fileStore,
	}
}

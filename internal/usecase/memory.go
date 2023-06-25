package usecase

type MemoryService struct {
	cache     CacheURL
	fileStore FileStoreURL
}

func NewMemoryService(cache CacheURL, fileStore FileStoreURL) *MemoryService {
	return &MemoryService{
		cache:     cache,
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
		_, err = m.cache.WriteURL(url, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

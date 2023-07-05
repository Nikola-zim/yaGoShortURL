package entity

type ConfigInit struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	UnitTestFlag    bool
	PostgresURL     string `env:"DATABASE_DSN"`
	UsingDB         bool
	DelBatch        int64
}

package static

type JSONApi struct {
	URL string
}

type JSONRes struct {
	Res string `json:"result"`
}

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

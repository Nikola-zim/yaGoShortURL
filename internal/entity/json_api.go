package entity

type JSONApi struct {
	URL string
}

type JSONRes struct {
	Res string `json:"result"`
}

type JSONAllInfo struct {
	ShortURL string `json:"short_url"`
	FullURL  string `json:"original_url"`
}

// BatchAPI - JSON API.
type BatchAPI struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResBatchAPI struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

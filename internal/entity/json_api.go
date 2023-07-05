package entity

type InputJSON struct {
	URL string
}

type ResultJSON struct {
	Res string `json:"result"`
}

type JSONAllInfo struct {
	ShortURL  string `json:"short_url"`
	FullURL   string `json:"original_url"`
	IsDeleted bool   `json:"is_deleted"`
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

type DeleteList struct {
	List []string
}

type DeleteMsg struct {
	List   []string
	UserID uint64
}

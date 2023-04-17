package static

type JSONApi struct {
	URL string
}

type JSONRes struct {
	Res string `json:"result"`
}

type JSONAllInfo struct {
	BaseURL string `json:"short_url"`
	FullURL string `json:"original_url"`
}

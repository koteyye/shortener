package models

type LongURL struct {
	URL string `json:"url"`
}

type ShortURL struct {
	Result string `json:"result"`
}

type FileString struct {
	ID          int    `json:"id"`
	ShortURL    string `json:"shortURL"`
	OriginalURL string `json:"originalURL"`
}

type URLList struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
}

type OriginURLList struct {
	ID        string `json:"correlation_id"`
	OriginURL string `json:"original_url"`
}

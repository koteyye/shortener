package models

type LongURL struct {
	URL string `json:"url"`
}

type ShortURL struct {
	Result string `json:"result"`
}

type FileString struct {
	Id          int    `json:"id"`
	ShortURL    string `json:"shortURL"`
	OriginalURL string `json:"originalURL"`
}

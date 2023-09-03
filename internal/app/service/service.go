package service

import "github.com/koteyye/shortener/internal/app/storage"

type Shortener interface {
	ShortURL(url string) (string, error)
	LongURL(shortURL string) (string, error)
}

type Service struct {
	Shortener
}

func NewService(storage *storage.UrlHandler) *Service {
	return &Service{
		Shortener: NewShortenerService(storage.UrlStorage),
	}
}

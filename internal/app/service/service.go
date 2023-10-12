package service

import (
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
)

type Shortener interface {
	ShortURL(url string) (string, error)
	LongURL(shortURL string) (string, error)
	Ping() error
}

type Service struct {
	Shortener
}

func NewService(storage *storage.URLHandler, shortener *config.Shortener) *Service {
	return &Service{
		Shortener: NewShortenerService(storage.URLStorage, shortener),
	}
}

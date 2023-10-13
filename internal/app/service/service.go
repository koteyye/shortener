package service

import (
	"context"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/storage"
)

type Shortener interface {
	ShortURL(ctx context.Context, url string) (string, error)
	LongURL(ctx context.Context, shortURL string) (string, error)
	Ping(ctx context.Context) error
	Batch(ctx context.Context, originalList []*models.OriginURLList) ([]*models.URLList, error)
}

type Service struct {
	Shortener
}

func NewService(storage *storage.URLHandler, shortener *config.Shortener) *Service {
	return &Service{
		Shortener: NewShortenerService(storage.URLStorage, shortener),
	}
}

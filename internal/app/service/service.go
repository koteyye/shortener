package service

import (
	"context"

	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/storage"
	"go.uber.org/zap"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Shortener interface {
	AddShortURL(ctx context.Context, url string, userID string) (string, error)
	GetOriginURL(ctx context.Context, shortURL string) (string, error)
	PingDB(ctx context.Context) error
	Batch(ctx context.Context, originalList []*models.OriginURLList, userID string) ([]*models.URLList, error)
	GetShortURLFromOriginal(ctx context.Context, originalURL string) (string, error)
	GetURLByUser(ctx context.Context, userID string) ([]*models.AllURLs, error)
	DeleteURLByUser(ctx context.Context, urls []string, userID string)
}

type Service struct {
	Shortener
}

func NewService(storage *storage.URLHandler, shortener *config.Shortener, logger *zap.SugaredLogger) *Service {
	return &Service{
		Shortener: NewShortenerService(storage.URLStorage, shortener, logger),
	}
}

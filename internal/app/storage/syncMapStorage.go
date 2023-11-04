package storage

import (
	"context"
	"github.com/koteyye/shortener/internal/app/models"
	"sync"
)

type URLMap struct {
	storage sync.Map
}

func NewURLMap() *URLMap {
	return &URLMap{storage: sync.Map{}}
}

func (u *URLMap) GetURLByUser(_ context.Context, _ string) ([]*models.AllURLs, error) {
	return nil, models.ErrMockNotSupported
}

func (u *URLMap) Ping(_ context.Context) error {
	return models.ErrMockNotSupported
}

func (u *URLMap) GetShortURL(_ context.Context, _ string) (string, error) {
	return "", models.ErrMockNotSupported
}

func (u *URLMap) AddURL(_ context.Context, k, s string, _ string) error {
	u.storage.Store(k, s)
	return nil
}

func (u *URLMap) GetURL(_ context.Context, k string) (string, error) {

	url, ok := u.storage.Load(k)
	if !ok {
		return "", models.ErrNotFound
	}
	return url.(string), nil

}

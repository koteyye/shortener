package storage

import (
	"context"
	"sync"

	"github.com/koteyye/shortener/internal/app/models"
)

// URLMap структура мок-хранилища.
type URLMap struct {
	storage sync.Map
}

// NewURLMap возвращает новый экземпляр мок-хранилища.
func NewURLMap() *URLMap {
	return &URLMap{storage: sync.Map{}}
}

// GetURLByUser возвращает список URL по текущему пользователю (не поддерживается).
func (u *URLMap) GetURLByUser(_ context.Context, _ string) ([]*models.URLList, error) {
	return nil, models.ErrMockNotSupported
}

// DeleteURLByUser удаляет сокращенный URL из поступающего канала (не поддерживается).
func (u *URLMap) DeleteURLByUser(_ context.Context, _ []string) error {
	return models.ErrMockNotSupported
}

// GetDBPing проверяет подключение к БД (не поддерживается).
func (u *URLMap) GetDBPing(_ context.Context) error {
	return models.ErrMockNotSupported
}

// GetShortURL получить сокращенный URL (не поддерживается).
func (u *URLMap) GetShortURL(_ context.Context, _ string) (string, error) {
	return "", models.ErrMockNotSupported
}

// AddURL добавляет сокращенный URL в хранилище.
func (u *URLMap) AddURL(_ context.Context, k, s string, _ string) error {
	u.storage.Store(k, s)
	return nil
}

// GetURL получить сокращенный URL из хранилища.
func (u *URLMap) GetURL(_ context.Context, k string) (*models.SingleURL, error) {

	url, ok := u.storage.Load(k)
	if !ok {
		return &models.SingleURL{}, models.ErrNotFound
	}
	return &models.SingleURL{URL: url.(string)}, nil
}

// GetCount получить количество URL и пользователей (не поддерживается)
func (u *URLMap) GetCount(_ context.Context) (int, int, error) {
	return 0, 0, models.ErrMockNotSupported
}

// BatchAddURL множественное добавление сокращенных URL в хранилище.
func (u *URLMap) BatchAddURL(_ context.Context, urlList []*models.URLList, _ string) error {
	for i := range urlList {
		u.storage.Store(urlList[i].ShortURL, urlList[i].URL)
	}
	return nil
}

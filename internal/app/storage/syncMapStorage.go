package storage

import (
	"context"
	"errors"
	"sync"
)

type URLMap struct {
	storage sync.Map
}

func NewURLMap() *URLMap {
	return &URLMap{storage: sync.Map{}}
}

func (u *URLMap) Ping(_ context.Context) error {
	return errors.New("не поддерживается на моках")
}

func (u *URLMap) GetShortURL(_ context.Context, _ string) (string, error) {
	return "", errors.New("не поддерживается на моках")
}

func (u *URLMap) AddURL(_ context.Context, k, s string) error {
	u.storage.Store(k, s)
	return nil
}

func (u *URLMap) GetURL(_ context.Context, k string) (string, error) {

	url, ok := u.storage.Load(k)
	if !ok {
		return "", ErrNotFound
	}
	return url.(string), nil

}

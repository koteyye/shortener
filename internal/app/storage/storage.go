package storage

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("не найдено такого значения")

type URLStorage interface {
	AddURL(string, string)
	GetURL(string) (string, error)
}

type URLHandler struct {
	URLStorage
}

type URLMap struct {
	storage sync.Map
}

func NewURLHandle() *URLHandler {
	return &URLHandler{
		URLStorage: &URLMap{
			storage: sync.Map{},
		},
	}
}

func (u *URLMap) AddURL(k, s string) {
	u.storage.Store(k, s)
}

func (u *URLMap) GetURL(k string) (string, error) {
	url, ok := u.storage.Load(k)
	if !ok {
		return "", ErrNotFound
	}
	return url.(string), nil
}

package storage

import (
	"errors"
	"sync"
)

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
	urlHandler := &URLHandler{
		&URLMap{
			storage: sync.Map{},
		},
	}
	return urlHandler
}

func (u *URLMap) AddURL(k, s string) {
	u.storage.Store(k, s)
}

func (u *URLMap) GetURL(k string) (string, error) {
	url, ok := u.storage.Load(k)
	if !ok {
		return "", errors.New("нет такого значения")
	}
	return url.(string), nil
}

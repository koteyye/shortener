package storage

import "errors"

type URLStorage interface {
	AddURL(string, string) (bool, error)
	GetURL(string) (string, error)
}

type URLHandler struct {
	URLStorage
}

type URLMap struct {
	storage map[string]string
}

func NewURLHandle() *URLHandler {
	urlHandler := &URLHandler{
		&URLMap{
			storage: make(map[string]string),
		},
	}
	return urlHandler
}

func (u URLMap) AddURL(k, s string) (bool, error) {
	u.storage[k] = s
	return true, nil
}

func (u URLMap) GetURL(k string) (string, error) {
	url, ok := u.storage[k]
	if !ok {
		return "", errors.New("нет такого значения")
	}
	return url, nil
}

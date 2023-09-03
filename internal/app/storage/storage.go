package storage

import "errors"

type UrlStorage interface {
	AddUrl(string, string) (bool, error)
	GetUrl(string) (string, error)
}

type UrlHandler struct {
	UrlStorage
}

type UrlMap struct {
	storage map[string]string
}

func NewUrlHandler() *UrlHandler {
	urlHandler := &UrlHandler{
		&UrlMap{
			storage: make(map[string]string),
		},
	}
	return urlHandler
}

func (u UrlMap) AddUrl(k, s string) (bool, error) {
	u.storage[k] = s
	return true, nil
}

func (u UrlMap) GetUrl(k string) (string, error) {
	url, ok := u.storage[k]
	if !ok {
		return "", errors.New("Нет такого значения")
	}
	return url, nil
}

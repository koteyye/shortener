package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/koteyye/shortener/internal/app/storage"
	"time"
)

type ShortenerService struct {
	storage storage.URLStorage
}

func NewShortenerService(storage storage.URLStorage) *ShortenerService {
	return &ShortenerService{storage: storage}
}

func (s ShortenerService) LongURL(shortURL string) (string, error) {
	res, err := s.storage.GetURL(shortURL)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s ShortenerService) ShortURL(url string) (string, error) {
	res := generateUnitKey()
	ok, _ := s.storage.AddURL(res, url)
	if !ok {
		return "", errors.New("не удалось запись значение в хранилище")
	}
	urlRes := "http://localhost:8080/" + res
	return urlRes, nil
}

func generateUnitKey() string {
	t := time.Now().UnixNano()

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(t)))
}

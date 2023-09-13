package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/koteyye/shortener/internal/app/storage"
	"regexp"
	"time"
)

var (
	ErrNullRequestBody       = errors.New("не указана ссылка для сокращения")
	ErrInvalidRequestBodyURL = errors.New("некорректно указана ссылка в запросе")
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
	if url == "" {
		return "", ErrNullRequestBody
	}
	if partURL, _ := regexp.Match(`(http)`, []byte(url)); !partURL {
		return "", ErrInvalidRequestBodyURL
	}

	res := generateUnitKey()
	s.storage.AddURL(res, url)
	urlRes := "http://localhost:8080/" + res
	return urlRes, nil
}

func generateUnitKey() string {
	t := time.Now().UnixNano()

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(t)))
}

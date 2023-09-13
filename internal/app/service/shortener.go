package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
	"regexp"
	"time"
)

var (
	ErrNullRequestBody       = errors.New("не указана ссылка для сокращения")
	ErrInvalidRequestBodyURL = errors.New("некорректно указана ссылка в запросе")
)

type ShortenerService struct {
	storage   storage.URLStorage
	shortener *config.Shortener
}

func NewShortenerService(storage storage.URLStorage, shortener *config.Shortener) *ShortenerService {
	return &ShortenerService{storage: storage, shortener: shortener}
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
	urlRes := s.shortener.Listen + s.shortener.BaseURL + res
	return urlRes, nil
}

func generateUnitKey() string {
	t := time.Now().UnixNano()

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(t)))
}

package service

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
	"net/url"
	"strings"
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

func (s ShortenerService) ShortURL(urlVal string) (string, error) {
	if urlVal == "" {
		return "", ErrNullRequestBody
	}
	if partURL := strings.Contains(urlVal, "http"); !partURL {
		return "", ErrInvalidRequestBodyURL
	}

	res := generateUnitKey()
	if err := s.storage.AddURL(res, urlVal); err != nil {
		return "", fmt.Errorf("add url: %w", err)
	}
	urlRes, err := url.JoinPath(s.shortener.Listen, "/", res)
	if err != nil {
		return "", err
	}
	return urlRes, nil
}

func generateUnitKey() string {
	t := time.Now().UnixNano()

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(t)))
}

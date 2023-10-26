package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/storage"
	"net/url"
	"time"
)

type ShortenerService struct {
	storage   storage.URLStorage
	shortener *config.Shortener
}

func NewShortenerService(storage storage.URLStorage, shortener *config.Shortener) *ShortenerService {
	return &ShortenerService{storage: storage, shortener: shortener}
}

func (s ShortenerService) Batch(ctx context.Context, originalList []*models.OriginURLList, userID string) ([]*models.URLList, error) {
	var urllist []*models.URLList
	for _, origin := range originalList {
		short, err := s.AddShortURL(ctx, origin.OriginURL, userID)
		if err != nil {
			if models.MapConflict(err) {
				shortURL, err := s.GetShortURLFromOriginal(ctx, origin.OriginURL)
				if err != nil {
					return nil, fmt.Errorf("ошибка при получении задублированного url: %v", err)
				}
				urllist = append(urllist, &models.URLList{ID: origin.ID, ShortURL: shortURL, Msg: models.ErrDuplicate.Error()})
				continue
			}
			return nil, fmt.Errorf("ошибка при заполнении url list: %v", err)
		}
		urllist = append(urllist, &models.URLList{ID: origin.ID, ShortURL: short})
	}
	return urllist, nil
}

func (s ShortenerService) PingDB(ctx context.Context) error {
	return s.storage.Ping(ctx)
}

func (s ShortenerService) GetOriginURL(ctx context.Context, shortURL string) (string, error) {
	res, err := s.storage.GetURL(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (s ShortenerService) AddShortURL(ctx context.Context, urlVal string, userID string) (string, error) {
	res := generateUnitKey()
	if err := s.storage.AddURL(ctx, res, urlVal, userID); err != nil {
		return "", fmt.Errorf("add url: %w", err)
	}
	urlRes, err := url.JoinPath(s.shortener.Listen, "/", res)
	if err != nil {
		return "", err
	}
	return urlRes, nil
}

func (s ShortenerService) GetShortURLFromOriginal(ctx context.Context, urlVal string) (string, error) {
	short, err := s.storage.GetShortURL(ctx, urlVal)
	if err != nil {
		return "", err
	}
	urlRes, err := url.JoinPath(s.shortener.Listen, "/", short)
	if err != nil {
		return "", err
	}
	return urlRes, nil
}

func (s ShortenerService) GetURLByUser(ctx context.Context, userID string) ([]*models.AllURLs, error) {
	allURLs, err := s.storage.GetURLByUser(ctx, userID)
	for _, urlItem := range allURLs {
		finalURL, err := url.JoinPath(s.shortener.Listen, "/", urlItem.ShortURL)
		if err != nil {
			return nil, err
		}
		urlItem.ShortURL = finalURL
	}
	return allURLs, err
}

func generateUnitKey() string {
	t := time.Now().UnixNano()

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(t)))
}

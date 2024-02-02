package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/koteyye/shortener/internal/app/models"
)

// Batch сокращение множества URL.
func (s Service) Batch(ctx context.Context, originalList []*models.URLList, userID string) ([]*models.URLList, error) {
	var urllist []*models.URLList
	for _, origin := range originalList {
		short, err := s.AddShortURL(ctx, origin.URL, userID)
		if err != nil {
			if models.MapConflict(err) {
				shortURL, err := s.GetShortURLFromOriginal(ctx, origin.URL)
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

// GetDBPing проверка подключения к БД.
func (s Service) GetDBPing(ctx context.Context) error {
	return s.storage.GetDBPing(ctx)
}

// GetOriginURL получение оригинального URL.
func (s Service) GetOriginURL(ctx context.Context, shortURL string) (string, error) {
	res, err := s.storage.GetURL(ctx, shortURL)
	if err != nil {
		return "", err
	}
	if res.IsDeleted {
		return "", models.ErrDeleted
	}
	return res.URL, nil
}

// AddShortURL сокращение оригинального URL.
func (s Service) AddShortURL(ctx context.Context, urlVal string, userID string) (string, error) {
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

// GetShortURLFromOriginal получение сокращенного URL по оригинальному.
func (s Service) GetShortURLFromOriginal(ctx context.Context, urlVal string) (string, error) {
	short, err := s.storage.GetShortURL(ctx, urlVal)
	if err != nil {
		return "", err
	}
	urlRes, err := url.JoinPath(s.shortener.Listen, short)
	if err != nil {
		return "", err
	}
	return urlRes, nil
}

// GetURLByUser получение списка URL по текущему пользователю.
func (s Service) GetURLByUser(ctx context.Context, userID string) ([]*models.URLList, error) {
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

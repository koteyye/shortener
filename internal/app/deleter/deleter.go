package deleter

import (
	"context"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/koteyye/shortener/internal/app/storage"
)

const batchMaxLen = 50

// Deleter воркер выполняющий асинхронное удаление
type Deleter struct {
	delURLch chan DeleteURL
	ticker   *time.Ticker
	storage  storage.URLStorage
	logger   *zap.SugaredLogger
	mutex    sync.Mutex
}

type DeleteURL struct {
	URL    []string
	UserID string
}

// StartDeleter запускает воркер
func InitDeleter(delURLch chan DeleteURL, storage storage.URLStorage, logger *zap.SugaredLogger) *Deleter {
	return &Deleter{delURLch: delURLch, storage: storage, logger: logger, mutex: sync.Mutex{}, ticker: time.NewTicker(time.Second * 3)}
}

// StartWorker запускает обработчик удаления URL
func (d *Deleter) StartWorker(ctx context.Context) {
	var batch []string
	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				d.storage.DeleteURLByUser(ctx, batch)
			}
			return
		case url := <-d.delURLch:
			validURL, err := d.validateURL(url.URL, url.UserID)
			if err != nil {
				d.logger.Infof("can't get urls by userID: %v, err: %w", url.UserID, err)
			}
			batch = append(batch, validURL...)
			if len(batch) >= batchMaxLen {
				d.storage.DeleteURLByUser(ctx, batch)
			}
		case <-d.ticker.C:
			if len(batch) > 0 {
				d.storage.DeleteURLByUser(ctx, batch)
			}

		}
	}
}

func (d *Deleter) validateURL(inURLs []string, userID string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var validURL []string
	urls, err := d.storage.GetURLByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	for idx := range inURLs {
		for idxURLItem := range urls {
			if strings.Contains(urls[idxURLItem].ShortURL, inURLs[idx]) {
				validURL = append(validURL, inURLs[idx])
			}
		}
	}
	return validURL, nil
}

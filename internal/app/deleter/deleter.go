package deleter

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/storage"
	"go.uber.org/zap"
)

const (
	maxURL = 50 // максимальное количество обрабатываемых URL

	maxItemMsg = "delete with maxitem"
	gracefulMsg = "stoped delete worker"
	timeoutMsg = "delete with timeout"
)

// Deleter воркер выполняющий асинхронное удаление
type Deleter struct {
	URL     chan string
	ticker  *time.Ticker
	storage storage.URLStorage
	logger  *zap.SugaredLogger
	mutex   sync.Mutex
	test *unitTest
}

type unitTest struct {
	isTest bool
	msg chan string
}

// StartDeleter запускает воркер
func InitDeleter(storage storage.URLStorage, logger *zap.SugaredLogger) *Deleter {
	return &Deleter{URL: make(chan string, maxURL), storage: storage, logger: logger, mutex: sync.Mutex{}, ticker: time.NewTicker(time.Second * 10), test: &unitTest{isTest: false}}
}

// Receive принимает URL в обработку
func (d *Deleter) Receive(ctx context.Context, delURLS []string, userID string) {
	urls, err := d.storage.GetURLByUser(ctx, userID)
	if err != nil {
		d.logger.Errorf("can't get urls by userID: %s, err: %w", userID, err)
		return
	}
	d.validateURL(ctx, delURLS, urls) // TODO another context must be
}

func (d *Deleter) validateURL(ctx context.Context, delURLS []string, urls []*models.URLList) {
	for idx := range delURLS {
		for idxUrlItem := range urls {
			if strings.Contains(urls[idxUrlItem].ShortURL, delURLS[idx]) {
				d.mutex.Lock()
				if len(d.URL) == 0 {
					d.ticker.Reset(time.Second * 10)
				}
				d.URL <- delURLS[idx]
				if len(d.URL) == maxURL {
					d.logger.Info("deleting url because of full capacity")
					urls := make([]string, 0, len(d.URL))
					for i := 0; i < len(d.URL); i++ {
						urls = append(urls, <-d.URL)
					}
					d.storage.DeleteURLByUser(ctx, urls)
				}
				d.mutex.Unlock()
				break
			}
		}
	}
}

func (d *Deleter) StartWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			{
				d.mutex.Lock()
				d.logger.Info("stopping deleter worker")
				urls := make([]string, 0, len(d.URL))
				for i := 0; i < len(d.URL); i++ {
					urls = append(urls, <-d.URL)
				}
				if d.test.isTest {
					d.test.msg <-gracefulMsg
					close(d.test.msg)
				}
				d.storage.DeleteURLByUser(ctx, urls)
				close(d.URL)
				d.ticker.Stop()
				d.mutex.Unlock()
				d.logger.Info("stopped deleter worker")

				return
			}
		case <-d.ticker.C:
			{
				if len(d.URL) != 0 {
					d.mutex.Lock()
					d.logger.Info("deleting url because of timer")
					urls := make([]string, 0, len(d.URL))
					for i := 0; i < len(d.URL); i++ {
						urls = append(urls, <-d.URL)
					}
					d.mutex.Unlock()
					if d.test.isTest {
						d.test.msg <-timeoutMsg
					}
					d.storage.DeleteURLByUser(ctx, urls)
				}
			}

		}
	}
}
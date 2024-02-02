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
)

// Deleter воркер выполняющий асинхронное удаление
type Deleter struct {
	URL chan string
	once sync.Once
	Counter *Counter
	storage storage.URLStorage
	logger *zap.SugaredLogger
	termCh chan struct{}
	wg *sync.WaitGroup
}

// StartDeleter запускает воркер
func StartDeleter(storage storage.URLStorage, logger *zap.SugaredLogger) *Deleter {
	urls := make(chan string)
	cnt := &Counter{
		num: 0,
	}
	return &Deleter{URL: urls, Counter: cnt, storage: storage, logger: logger, termCh: make(chan struct{}), wg: &sync.WaitGroup{}}
}

// Receive принимает URL в обработку
func (d *Deleter) Receive(delURLS []string, userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if d.Counter.Value() > maxURL {
		// хз как организовать ожидание 
	}
	d.Counter.Inc(len(delURLS))

	urls, err := d.storage.GetURLByUser(ctx, userID)
	if err != nil {
		d.logger.Errorf("can't get urls by userID: %s, err: %w", userID, err)
		d.Counter.Dec(len(delURLS))
		return
	}
	d.validateURL(delURLS, urls)
}

func (d *Deleter) validateURL(delURLS []string, urls []*models.URLList) {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for _, delURL := range delURLS {
			var validURL string
			for _, urlItem := range urls {
				if strings.Contains(urlItem.ShortURL, delURL) {
					validURL = delURL
					break
				}
			}

			select {
			case <- d.termCh:
				return
			case d.URL <- validURL:
			}
		}
	}()
	d.wg.Wait()
}

func (d *Deleter) execute() {
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		for data := range d.URL {
			var batchDeleteURL []string
			timer := time.NewTimer(30 * time.Second)

			batchDeleteURL= append(batchDeleteURL, data)
			<- timer.C

			d.storage.DeleteURLByUser(context.Background(), batchDeleteURL)
		}
	}()
	
	go func() {
		d.wg.Wait()
	}()
}

func (d *Deleter) Close() {
	if d.Counter.Value() != 0 {
		// как сделать ожидание 0 ?
	}
}
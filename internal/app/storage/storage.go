package storage

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/shortener/internal/app/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock.go
// URLStorage интерфейс хранилища.
type URLStorage interface {
	AddURL(context.Context, string, string, string) error
	GetURL(context.Context, string) (*models.URL, error)
	Ping(ctx context.Context) error
	GetShortURL(context.Context, string) (string, error)
	GetURLByUser(context.Context, string) ([]*models.AllURLs, error)
	DeleteURLByUser(context.Context, chan string) error
}

// URLHandler структура обработчика URL
type URLHandler struct {
	URLStorage
}

// NewURLHandle возвращает новый экземпляр обработчика URL
func NewURLHandle(log *zap.SugaredLogger, db *sqlx.DB, filePath string) (*URLHandler, error) {
	if db != nil {
		log.Info("start storage in db")
		return &URLHandler{
			URLStorage: NewPostgres(db),
		}, nil
	} else if filePath != "" {
		log.Info(fmt.Sprintf("start storage in file: %v", filePath))
		return &URLHandler{
			URLStorage: NewFileStorage(filePath),
		}, nil
	}
	log.Info("start storage in moc")
	return &URLHandler{
		URLStorage: NewURLMap(),
	}, nil
}

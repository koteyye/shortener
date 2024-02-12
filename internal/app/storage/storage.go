package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"github.com/koteyye/shortener/internal/app/models"
)

// URLStorage интерфейс хранилища.
//
//go:generate mockgen -source=storage.go -destination=mocks/mock.go
type URLStorage interface {
	AddURL(context.Context, string, string, string) error
	GetURL(context.Context, string) (*models.SingleURL, error)
	GetDBPing(ctx context.Context) error
	GetShortURL(context.Context, string) (string, error)
	GetURLByUser(context.Context, string) ([]*models.URLList, error)
	DeleteURLByUser(context.Context, []string) error
	GetCount(context.Context) (int, int, error)
	BatchAddURL(context.Context, []*models.URLList, string) error
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

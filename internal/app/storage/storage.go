package storage

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/shortener/internal/app/models"
	"go.uber.org/zap"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock.go

type URLStorage interface {
	AddURL(context.Context, string, string) error
	GetURL(context.Context, string) (string, error)
	Ping(ctx context.Context) error
	GetShortURL(context.Context, string) (string, error)
	GetURLByUser(context.Context, string) ([]*models.AllURLs, error)
}

type URLHandler struct {
	URLStorage
}

func NewURLHandle(ctx context.Context, db *sqlx.DB, filePath string) (*URLHandler, error) {
	log := ctx.Value("logger").(zap.SugaredLogger)
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

package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var ErrNotFound = errors.New("не найдено такого значения")

type URLStorage interface {
	AddURL(context.Context, string, string) error
	GetURL(context.Context, string) (string, error)
	Ping(ctx context.Context) error
	GetShortURL(context.Context, string) (string, error)
}

type URLHandler struct {
	URLStorage
}

func NewURLHandle(db *sqlx.DB, filePath string) (*URLHandler, error) {
	newLogger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	defer newLogger.Sync()

	logger := *newLogger.Sugar()
	if db != nil {
		logger.Info("start storage in db")
		return &URLHandler{
			URLStorage: NewPostgres(db),
		}, nil
	} else if filePath != "" {
		logger.Info(fmt.Sprintf("start storage in file: %v", filePath))
		return &URLHandler{
			URLStorage: NewFileStorage(filePath),
		}, nil
	}
	logger.Info("start storage in moc")
	return &URLHandler{
		URLStorage: NewURLMap(),
	}, nil
}

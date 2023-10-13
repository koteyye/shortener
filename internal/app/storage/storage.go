package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"sync"
)

var ErrNotFound = errors.New("не найдено такого значения")

type URLStorage interface {
	AddURL(context.Context, string, string) error
	GetURL(context.Context, string) (string, error)
	Ping(ctx context.Context) error
}

type URLHandler struct {
	URLStorage
}

type URLMap struct {
	storage     sync.Map
	fileStorage *FileStorage
}

func NewURLHandle(db *sqlx.DB, filePath string) *URLHandler {
	newLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer newLogger.Sync()

	logger := *newLogger.Sugar()
	if db != nil {
		logger.Info("start storage in db")
		return &URLHandler{
			URLStorage: NewPostgres(db),
		}
	} else if filePath != "" {
		logger.Info(fmt.Sprintf("start storage in file: %v", filePath))
		return &URLHandler{
			URLStorage: NewFileStorage(filePath),
		}
	}
	logger.Info("start storage in moc")
	return &URLHandler{
		URLStorage: &URLMap{
			storage: sync.Map{},
		},
	}
}

func (u *URLMap) Ping(_ context.Context) error {
	return errors.New("в качестве бд используется мок")
}

func (u *URLMap) AddURL(_ context.Context, k, s string) error {
	u.storage.Store(k, s)
	return nil
}

func (u *URLMap) GetURL(_ context.Context, k string) (string, error) {

	url, ok := u.storage.Load(k)
	if !ok {
		return "", ErrNotFound
	}
	return url.(string), nil

}

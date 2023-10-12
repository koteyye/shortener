package storage

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"sync"
)

var ErrNotFound = errors.New("не найдено такого значения")

type URLStorage interface {
	AddURL(string, string) error
	GetURL(string) (string, error)
	Ping() error
}

type URLHandler struct {
	URLStorage
}

type URLMap struct {
	storage     sync.Map
	fileStorage *FileStorage
}

func (u *URLMap) Ping() error {
	return errors.New("в качестве бд используется мок")
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

func (u *URLMap) AddURL(k, s string) error {
	u.storage.Store(k, s)
	return nil
}

func (u *URLMap) GetURL(k string) (string, error) {

	url, ok := u.storage.Load(k)
	if !ok {
		return "", ErrNotFound
	}
	return url.(string), nil

}

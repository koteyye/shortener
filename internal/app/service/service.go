package service

import (
	"go.uber.org/zap"

	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
)

// Serivce структура сервисного слоя
type Service struct {
	storage   storage.URLStorage
	shortener *config.Shortener
	logger    *zap.SugaredLogger
}

// NewService возвращает новый экземпляр Service
func NewService(storage storage.URLStorage, shortener *config.Shortener, logger *zap.SugaredLogger) *Service {
	return &Service{storage: storage, shortener: shortener, logger: logger}
}

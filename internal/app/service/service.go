package service

import (
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
	"go.uber.org/zap"
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

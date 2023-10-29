package models

import (
	"errors"
	"github.com/lib/pq"
)

const (
	PqDuplicateErr  = "23505"
	JWTExpiredToken = 16
)

var (
	ErrNullRequestBody       = errors.New("в запросе нет сокращенной ссылки")
	ErrInvalidRequestBodyURL = errors.New("некорректно указана ссылка в запросе")
	ErrNotFound              = errors.New("не найдено такого значения")
	ErrDuplicate             = errors.New("в бд уже есть сокращенный url")
	ErrDB                    = errors.New("непредвиденная ошибка в бд")
	ErrInvalidToken          = errors.New("невалидный токен")
	ErrExpiredToken          = errors.New("токен просрочен")
	ErrMockNotSupported      = errors.New("не поддерживается в моках")
	ErrFileNotSupported      = errors.New("не поддерживается в файле")
	ErrDeleted               = errors.New("ссылка удалена")
)

func MapConflict(err error) bool {
	var errPQ *pq.Error
	if errors.As(err, &errPQ) {
		if errPQ.Code == PqDuplicateErr {
			return true
		}
	}
	return false
}

func MapBatchConflict(list []*URLList) bool {
	msgCount := 0
	for _, item := range list {
		if item.Msg != "" {
			msgCount += 1
		}
	}
	return msgCount != len(list)
}

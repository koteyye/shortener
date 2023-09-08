package service

import (
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Тестовый конфиг для сокращения ссылок
var shortenerCfg = &config.Shortener{
	BaseURL: "/",
	Listen:  "http://localhost:8080",
}

func TestShortenerService_LongURL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "positive",
			value: "bfebrehbrehgbre",
			want:  "https://practicum.yandex.ru/",
		},
		{
			name:  "negative",
			value: "fpowdkfgpowedr",
			want:  "",
		},
	}

	storages := storage.NewURLHandle()
	s := NewService(storages, shortenerCfg)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Добавляем в маппу тестовую ссылку
			storages.AddURL("bfebrehbrehgbre", "https://practicum.yandex.ru/")
			url, err := s.LongURL(test.value)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.want, url)
			}
		})
	}
}

func TestShortenerService_ShortURL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "positive",
			value: "https://practicum.yandex.ru/",
			want:  "http://localhost:8080/",
		},
		{
			name:  "null var",
			value: "",
			want:  "",
		},
		{
			name:  "invalid URL",
			value: "practicum.yandex.ru/",
			want:  "",
		},
	}
	storages := storage.NewURLHandle()
	s := NewService(storages, shortenerCfg)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url, err := s.ShortURL(test.value)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Contains(t, url, test.want)
			}
		})
	}
}

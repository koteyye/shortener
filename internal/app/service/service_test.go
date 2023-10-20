package service

import (
	"context"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

// Тестовый конфиг для сокращения ссылок
var shortenerCfg = &config.Shortener{
	Listen: "http://localhost:8080/",
}

func TestShortenerService_LongURL(t *testing.T) {
	testDir := t.TempDir()
	file, err := os.CreateTemp(testDir, "db")
	require.NoError(t, err)
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

	storages, _ := storage.NewURLHandle(nil, file.Name())
	s := NewService(storages, shortenerCfg)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Добавляем в маппу тестовую ссылку
			storages.AddURL(context.Background(), "bfebrehbrehgbre", "https://practicum.yandex.ru/")
			url, err := s.LongURL(context.Background(), test.value)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, test.want, url)
			}
		})
	}
	assert.NoError(t, err)
}

func TestShortenerService_ShortURL(t *testing.T) {
	testDir := t.TempDir()
	file, err := os.CreateTemp(testDir, "db")
	require.NoError(t, err)
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
	storages, _ := storage.NewURLHandle(nil, file.Name())
	s := NewService(storages, shortenerCfg)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url, err := s.ShortURL(context.Background(), test.value)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Contains(t, url, test.want)
			}
		})
	}
	err = os.RemoveAll("tmp/")
	assert.NoError(t, err)
}

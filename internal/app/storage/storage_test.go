package storage

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestStorage_AddURL(t *testing.T) {

	tests := []struct {
		name string
		key  string
		val  string
		want string
	}{
		{
			name: "positive",
			key:  "bfebrehbrehgbre",
			val:  "https://practicum.yandex.ru/",
			want: "https://practicum.yandex.ru/",
		},
	}

	s := NewURLHandle(strings.TrimLeft("/tmp/short-url-db.json", "/"))
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.AddURL(test.key, test.val)

			result, err := s.GetURL(test.key)

			assert.NoError(t, err)
			assert.Equal(t, test.want, result)
		})
	}

	err := os.RemoveAll("tmp/")
	if err != nil {
		return
	}
}

func TestStorage_GetURL(t *testing.T) {
	s := NewURLHandle(strings.TrimLeft("/tmp/short-url-db.json", "/"))
	//Кладем значение для теста
	s.AddURL("sdvgdsgv", "https://practicum.yandex.ru/")

	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:    "positive",
			key:     "sdvgdsgv",
			want:    "https://practicum.yandex.ru/",
			wantErr: false,
		},
		{
			name:    "negative",
			key:     "aweasdw",
			want:    "",
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.GetURL(test.key)

			result, err := s.GetURL(test.key)

			if !test.wantErr {
				assert.Equal(t, test.want, result)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}

	err := os.RemoveAll("tmp/")
	if err != nil {
		return
	}
}

package storage

import (
	"github.com/stretchr/testify/assert"
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

	s := NewURLHandle(nil, "")
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s.AddURL(test.key, test.val)

			result, err := s.GetURL(test.key)

			assert.NoError(t, err)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestStorage_GetURL(t *testing.T) {
	s := NewURLHandle(nil, "")
	//Кладем значение для теста
	err := s.AddURL("sdvgdsgv", "https://practicum.yandex.ru/")
	assert.NoError(t, err)

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
			url, err := s.GetURL(test.key)

			if !test.wantErr {
				assert.Equal(t, test.want, url)
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

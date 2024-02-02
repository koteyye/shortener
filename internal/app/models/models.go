package models

import (
	"fmt"
	"strings"
)

// SignleURL структура для одного URL
type SingleURL struct {
	ID        string `json:"id,omitempty"`
	URL       string `json:"url,omitempty" db:"originalurl"`
	ShortURL  string `json:"short_url,omitempty" db:"shorturl"`
	IsDeleted bool   `json:"is_deleted,omitempty" db:"is_deleted"`
}

// URLList структура для списка URL
type URLList struct {
	Number   int    `json:"id,omitempty"`
	ID       string `json:"correlation_id,omitempty"`
	URL      string `json:"original_url,omitempty" db:"originalurl"`
	ShortURL string `json:"short_url,omitempty" db:"shorturl"`
	Msg      string `json:"msg,omitempty"`
}

// Validate валидация списка с оригинальным URL.
func (o *URLList) Validate() error {
	if o.ID == "" {
		return ErrNullRequestBody
	}
	if o.URL == "" {
		return ErrInvalidRequestBodyURL
	}
	if !strings.Contains(o.URL, "http") {
		return fmt.Errorf("некорректно указан original_url с correlation_id: %v", o.ID)
	}
	return nil
}

// Validate валидация оригинального URL.
func (u *SingleURL) Validate() error {
	if u.URL == "" {
		return ErrNullRequestBody
	}
	if !strings.Contains(u.URL, "http") {
		return ErrInvalidRequestBodyURL
	}
	return nil
}

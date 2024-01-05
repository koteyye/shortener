package models

import (
	"fmt"
	"strings"
)

// OriginalURL оригинальный URL.
type OriginalURL struct {
	URL string `json:"url"`
}

// ShortURL сокращенный URL.
type ShortURL struct {
	Result string `json:"result"`
}

// AllURLs элемент списка со всеми URL.
type AllURLs struct {
	ID          int    `json:"id,omitempty"`
	ShortURL    string `json:"short_url" db:"shorturl"`
	OriginalURL string `json:"original_url" db:"originalurl"`
}

// URLList элемент списка с сокращенным URL.
type URLList struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
	Msg      string `json:"msg,omitempty"`
}

// OriginURLList элемент списка с оригинальным URL.
type OriginURLList struct {
	ID        string `json:"correlation_id"`
	OriginURL string `json:"original_url"`
}

// URL элемент списка с URL и признаком удаления.
type URL struct {
	ShortURL    string `db:"shorturl"`
	OriginalURL string `db:"originalurl"`
	IsDeleted   bool   `db:"is_deleted"`
}

// Validate валидация списка с оригинальным URL.
func (o *OriginURLList) Validate() error {
	if o.ID == "" {
		return ErrNullRequestBody
	}
	if o.OriginURL == "" {
		return ErrInvalidRequestBodyURL
	}
	if !strings.Contains(o.OriginURL, "http") {
		return fmt.Errorf("некорректно указан original_url с correlation_id: %v", o.ID)
	}
	return nil
}

// Validate валидация оригинального URL.
func (u *OriginalURL) Validate() error {
	if u.URL == "" {
		return ErrNullRequestBody
	}
	if !strings.Contains(u.URL, "http") {
		return ErrInvalidRequestBodyURL
	}
	return nil
}

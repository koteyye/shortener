package models

import (
	"fmt"
	"strings"
)

type OriginalURL struct {
	URL string `json:"url"`
}

type ShortURL struct {
	Result string `json:"result"`
}

type AllURLs struct {
	ID          int    `json:"id,omitempty"`
	ShortURL    string `json:"short_url" db:"shorturl"`
	OriginalURL string `json:"original_url" db:"originalurl"`
}

type URLList struct {
	ID       string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
	Msg      string `json:"msg,omitempty"`
}

type OriginURLList struct {
	ID        string `json:"correlation_id"`
	OriginURL string `json:"original_url"`
}

type URL struct {
	ShortURL    string `db:"shorturl"`
	OriginalURL string `db:"originalurl"`
	IsDeleted   bool   `db:"is_deleted"`
}

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

func (u *OriginalURL) Validate() error {
	if u.URL == "" {
		return ErrNullRequestBody
	}
	if !strings.Contains(u.URL, "http") {
		return ErrInvalidRequestBodyURL
	}
	return nil
}

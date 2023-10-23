package models

import (
	"errors"
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

func (o *OriginURLList) Validate() error {
	if o.ID == "" {
		return ErrNullRequestBody
	}
	if o.OriginURL == "" {
		return ErrInvalidRequestBodyURL
	}
	if !strings.Contains(o.OriginURL, "http") {
		return errors.New(fmt.Sprintf("некорректно указан original_url с correlation_id: %s", o.ID))
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

package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/koteyye/shortener/internal/app/models"
)

func mapRequestShortenURL(r *http.Request) (string, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать запрос: %v", err)
	}
	strReqBody := buf.String()
	strReqBody = strings.TrimSpace(strReqBody)

	if strReqBody == "" {
		return "", models.ErrNullRequestBody
	}
	if partURL := strings.Contains(strReqBody, "http"); !partURL {
		return "", models.ErrInvalidRequestBodyURL
	}

	return strReqBody, nil
}

func mapRequestJSONShortenURL(r *http.Request) (*models.SingleURL, error) {
	var input *models.SingleURL
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("невозможно прочитать запрос: %v", err)
	}
	if err := json.Unmarshal(body, &input); err != nil {
		return nil, fmt.Errorf("невозможно десериализировать запрос: %v", err)
	}
	err = input.Validate()
	if err != nil {
		return nil, err
	}
	return input, nil
}

func mapRequestBatch(r *http.Request) ([]*models.URLList, error) {
	var input []*models.URLList
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("невозможно прочитать запрос: %v", err)
	}
	if err := json.Unmarshal(body, &input); err != nil {
		return nil, fmt.Errorf("невозможно десериализировать запрос: %v", err)
	}
	for _, item := range input {
		err := item.Validate()
		if err != nil {
			return nil, err
		}
		item.ID = strings.TrimSpace(item.ID)
		item.URL = strings.TrimSpace(item.URL)
	}
	return input, nil
}

func mapRequestDeleteByUser(r *http.Request) ([]string, error) {
	var urls []string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("невозможно прочитать запрос: %v", err)
	}
	if err := json.Unmarshal(body, &urls); err != nil {
		return nil, fmt.Errorf("невозможно десериализовать запрос: %v", err)
	}
	return urls, nil
}

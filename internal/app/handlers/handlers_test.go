package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/service"
	mock_service "github.com/koteyye/shortener/internal/app/service/mocks"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const (
	testSecretKey = "super_secret_key"
)

// testInitHandler инициализация тестового обработчика
func testInitHandler(t *testing.T) (*Handlers, *mock_service.MockShortener) {
	c := gomock.NewController(t)
	defer c.Finish()

	repo := mock_service.NewMockShortener(c)
	service := service.Service{Shortener: repo}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log := *logger.Sugar()

	handler := NewHandlers(&service, &log, testSecretKey)

	return handler, repo
}

func TestHandlers_Batch(t *testing.T) {
	testRequest := `[
		{
			"correlation_id": "afd90f2c-b0df-4873-8ded-62d8e99593ba",
			"original_url": "http://lcpjtoddpyyp.yandex/sjkhh"
		},
		{
			"correlation_id": "f5993975-38fc-4f95-8234-0639079194cf",
			"original_url": "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv"
		}
		]`

	t.Run("batch", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			h, s := testInitHandler(t)

			r := httptest.NewRequest(http.MethodPost, "/batch", strings.NewReader(testRequest))
			w := httptest.NewRecorder()

			var testOriginURLList []*models.OriginURLList
			err := json.Unmarshal([]byte(testRequest), &testOriginURLList)
			assert.NoError(t, err)

			userID := uuid.NewString()
			ctx := context.WithValue(r.Context(), userIDKey, userID)

			s.EXPECT().Batch(gomock.Any(), testOriginURLList, userID).Return([]*models.URLList{
				{
					ID:       "afd90f2c-b0df-4873-8ded-62d8e99593ba",
					ShortURL: "http://lcpjtoddpyyp.yandex/sjkhh",
				},
				{
					ID:       "5993975-38fc-4f95-8234-0639079194cf",
					ShortURL: "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv",
				},
			}, nil)

			h.Batch(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusCreated, w.Code)
		})
		t.Run("conflict", func(t *testing.T) {
			h, s := testInitHandler(t)

			r := httptest.NewRequest(http.MethodPost, "/batch", strings.NewReader(testRequest))
			w := httptest.NewRecorder()

			var testOriginURLList []*models.OriginURLList
			err := json.Unmarshal([]byte(testRequest), &testOriginURLList)
			assert.NoError(t, err)

			userID := uuid.NewString()
			ctx := context.WithValue(r.Context(), userIDKey, userID)

			s.EXPECT().Batch(gomock.Any(), testOriginURLList, userID).Return([]*models.URLList{
				{
					ID:       "afd90f2c-b0df-4873-8ded-62d8e99593ba",
					ShortURL: "http://lcpjtoddpyyp.yandex/sjkhh",
					Msg:      models.ErrDuplicate.Error(),
				},
				{
					ID:       "5993975-38fc-4f95-8234-0639079194cf",
					ShortURL: "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv",
					Msg:      models.ErrDuplicate.Error(),
				},
			}, nil)

			h.Batch(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusConflict, w.Code)
		})
		t.Run("url_empty", func(t *testing.T) {
			h, _ := testInitHandler(t)

			reqBody := `[
				{
					"original_url": "http://lcpjtoddpyyp.yandex/sjkhh"
				},
				{
					"correlation_id": "f5993975-38fc-4f95-8234-0639079194cf",
					"original_url": "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv"
				}
				]`

			r := httptest.NewRequest(http.MethodPost, "/batch", strings.NewReader(reqBody))
			w := httptest.NewRecorder()

			var testOriginURLList []*models.OriginURLList
			err := json.Unmarshal([]byte(testRequest), &testOriginURLList)
			assert.NoError(t, err)

			userID := uuid.NewString()
			ctx := context.WithValue(r.Context(), userIDKey, userID)

			h.Batch(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})
}


func TestHandlers_GetURLsByUser(t *testing.T) {
	type mockBehavior func(r *mock_service.MockShortener, userID string)
	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "success",
			mockBehavior: func(r *mock_service.MockShortener, userID string) {
				r.EXPECT().GetURLByUser(gomock.Any(), userID).Return([]*models.AllURLs{
					{
						ShortURL:    "http://localhost:8080/iwvwpvwepofpwoe",
						OriginalURL: "http://yandex.ru",
					},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"short_url":"http://localhost:8080/iwvwpvwepofpwoe","original_url":"http://yandex.ru"}]`,
		},
		{
			name: "no content",
			mockBehavior: func(r *mock_service.MockShortener, userID string) {
				r.EXPECT().GetURLByUser(gomock.Any(), userID).Return(nil, sql.ErrNoRows)
			},
			expectedStatusCode: 204,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())

			repo := mock_service.NewMockShortener(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, userID.String())
			}

			services := service.Service{Shortener: repo}
			handler := Handlers{services: &services}

			handler.GetURLsByUser(w, r.WithContext(ctx))

			assert.Equal(t, test.expectedStatusCode, w.Code)
			if test.expectedResponseBody != "" {
				assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestHandlers_JSONShortenURL(t *testing.T) {
	type mockBehavior func(r *mock_service.MockShortener, url string)
	tests := []struct {
		name                 string
		inputBody            io.Reader
		inputURL             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "success",
			inputBody: strings.NewReader(`{"url": "http://yandex.ru"}`),
			inputURL:  "http://yandex.ru",
			mockBehavior: func(r *mock_service.MockShortener, url string) {
				r.EXPECT().AddShortURL(gomock.Any(), url, gomock.Any()).Return("http://localhost:8080/3g2gf2f2", nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"result": "http://localhost:8080/3g2gf2f2"}`,
		},
		{
			name:      "conflict",
			inputBody: strings.NewReader(`{"url": "http://yandex.ru"}`),
			inputURL:  "http://yandex.ru",
			mockBehavior: func(r *mock_service.MockShortener, url string) {
				r.EXPECT().AddShortURL(gomock.Any(), url, gomock.Any()).Return("", &pq.Error{Code: models.PqDuplicateErr})
			},
			expectedStatusCode:   409,
			expectedResponseBody: `{"result": "http://localhost:8080/3g2gf2f2"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := httptest.NewRequest(http.MethodGet, "/api/shorten", test.inputBody)
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())

			repo := mock_service.NewMockShortener(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, test.inputURL)
			}
			if test.expectedStatusCode == 409 {
				repo.EXPECT().GetShortURLFromOriginal(gomock.Any(), test.inputURL).Return("http://localhost:8080/3g2gf2f2", nil)
			}

			services := service.Service{Shortener: repo}
			handler := Handlers{services: &services}

			handler.JSONShortenURL(w, r.WithContext(ctx))

			assert.Equal(t, test.expectedStatusCode, w.Code)
			if test.expectedResponseBody != "" {
				assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
			}
		})
	}
}

func TestHandlers_ShortenURL(t *testing.T) {
	type mockBehavior func(r *mock_service.MockShortener, url string)
	tests := []struct {
		name                 string
		inputBody            io.Reader
		inputURL             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "success",
			inputBody: strings.NewReader(`http://yandex.ru`),
			inputURL:  "http://yandex.ru",
			mockBehavior: func(r *mock_service.MockShortener, url string) {
				r.EXPECT().AddShortURL(gomock.Any(), url, gomock.Any()).Return("http://localhost:8080/3g2gf2f2", nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `http://localhost:8080/3g2gf2f2`,
		},
		{
			name:      "conflict",
			inputBody: strings.NewReader(`http://yandex.ru`),
			inputURL:  "http://yandex.ru",
			mockBehavior: func(r *mock_service.MockShortener, url string) {
				r.EXPECT().AddShortURL(gomock.Any(), url, gomock.Any()).Return("", &pq.Error{Code: models.PqDuplicateErr})
			},
			expectedStatusCode:   409,
			expectedResponseBody: `http://localhost:8080/3g2gf2f2`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := httptest.NewRequest(http.MethodPost, "/", test.inputBody)
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())

			repo := mock_service.NewMockShortener(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, test.inputURL)
			}
			if test.expectedStatusCode == 409 {
				repo.EXPECT().GetShortURLFromOriginal(gomock.Any(), test.inputURL).Return("http://localhost:8080/3g2gf2f2", nil)
			}

			services := service.Service{Shortener: repo}
			handler := Handlers{services: &services}

			handler.ShortenURL(w, r.WithContext(ctx))

			assert.Equal(t, test.expectedStatusCode, w.Code)
			if test.expectedResponseBody != "" {
				assert.Equal(t, test.expectedResponseBody, w.Body.String())
			}
		})
	}
}

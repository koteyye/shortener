package handlers

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/service"
	mock_service "github.com/koteyye/shortener/internal/app/service/mocks"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlers_Batch(t *testing.T) {
	type mockBehavior func(r *mock_service.MockShortener, urlList []*models.OriginURLList)
	tests := []struct {
		name                 string
		inputBody            io.Reader
		inputOriginURLList   []*models.OriginURLList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "success",
			inputBody: strings.NewReader(`[
				{
					"correlation_id": "afd90f2c-b0df-4873-8ded-62d8e99593ba",
					"original_url": "http://lcpjtoddpyyp.yandex/sjkhh"
				},
				{
					"correlation_id": "f5993975-38fc-4f95-8234-0639079194cf",
					"original_url": "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv"
				}
				]`),
			inputOriginURLList: []*models.OriginURLList{
				{
					ID:        "afd90f2c-b0df-4873-8ded-62d8e99593ba",
					OriginURL: "http://lcpjtoddpyyp.yandex/sjkhh",
				},
				{
					ID:        "f5993975-38fc-4f95-8234-0639079194cf",
					OriginURL: "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv",
				},
			},
			mockBehavior: func(r *mock_service.MockShortener, urlList []*models.OriginURLList) {
				r.EXPECT().Batch(gomock.Any(), urlList).Return([]*models.URLList{
					{
						ID:       "afd90f2c-b0df-4873-8ded-62d8e99593ba",
						ShortURL: "http://lcpjtoddpyyp.yandex/sjkhh",
					},
					{
						ID:       "5993975-38fc-4f95-8234-0639079194cf",
						ShortURL: "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv",
					},
				}, nil)
			},
			expectedStatusCode: 201,
			expectedResponseBody: `[
				{
					"correlation_id": "afd90f2c-b0df-4873-8ded-62d8e99593ba",
					"short_url": "http://lcpjtoddpyyp.yandex/sjkhh"
				},
				{
					"correlation_id": "5993975-38fc-4f95-8234-0639079194cf",
					"short_url": "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv"
				}
			]`,
		},
		{
			name: "conflict",
			inputBody: strings.NewReader(`[
				{
					"correlation_id": "afd90f2c-b0df-4873-8ded-62d8e99593ba",
					"original_url": "http://lcpjtoddpyyp.yandex/sjkhh"
				},
				{
					"correlation_id": "f5993975-38fc-4f95-8234-0639079194cf",
					"original_url": "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv"
				}
				]`),
			inputOriginURLList: []*models.OriginURLList{
				{
					ID:        "afd90f2c-b0df-4873-8ded-62d8e99593ba",
					OriginURL: "http://lcpjtoddpyyp.yandex/sjkhh",
				},
				{
					ID:        "f5993975-38fc-4f95-8234-0639079194cf",
					OriginURL: "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv",
				},
			},
			mockBehavior: func(r *mock_service.MockShortener, urlList []*models.OriginURLList) {
				r.EXPECT().Batch(gomock.Any(), urlList).Return([]*models.URLList{
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
			},
			expectedStatusCode: 409,
			expectedResponseBody: `[
				{
					"correlation_id": "afd90f2c-b0df-4873-8ded-62d8e99593ba",
					"short_url": "http://lcpjtoddpyyp.yandex/sjkhh",
					"msg": "в бд уже есть сокращенный url"
				},
				{
					"correlation_id": "5993975-38fc-4f95-8234-0639079194cf",
					"short_url": "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv",
					"msg": "в бд уже есть сокращенный url"
				}
			]`,
		},
		{
			name: "errNullBody",
			inputBody: strings.NewReader(`[
				{
					"original_url": "http://lcpjtoddpyyp.yandex/sjkhh"
				},
				{
					"correlation_id": "f5993975-38fc-4f95-8234-0639079194cf",
					"original_url": "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv"
				}
				]`),
			inputOriginURLList: []*models.OriginURLList{},
			expectedStatusCode: 400,
			expectedResponseBody: `{
			"Message": "в запросе нет сокращенной ссылки"
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := httptest.NewRequest(http.MethodPost, "/batch", test.inputBody)
			r.Header.Set("Content-Type", ctApplicationJSON)
			w := httptest.NewRecorder()

			repo := mock_service.NewMockShortener(c)
			if test.mockBehavior != nil {
				test.mockBehavior(repo, test.inputOriginURLList)
			}

			services := service.Service{Shortener: repo}
			handler := Handlers{services: &services}
			handler.Batch(w, r)

			result := w.Result()
			defer result.Body.Close()

			body, err := io.ReadAll(result.Body)
			assert.NoError(t, err)

			assert.Equal(t, test.expectedStatusCode, result.StatusCode)
			assert.JSONEq(t, test.expectedResponseBody, string(body))
		})
	}
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
				r.EXPECT().AddShortURL(gomock.Any(), url).Return("http://localhost:8080/3g2gf2f2", nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"result": "http://localhost:8080/3g2gf2f2"}`,
		},
		{
			name:      "conflict",
			inputBody: strings.NewReader(`{"url": "http://yandex.ru"}`),
			inputURL:  "http://yandex.ru",
			mockBehavior: func(r *mock_service.MockShortener, url string) {
				r.EXPECT().AddShortURL(gomock.Any(), url).Return("", &pq.Error{Code: models.PqDuplicateErr})
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
				r.EXPECT().AddShortURL(gomock.Any(), url).Return("http://localhost:8080/3g2gf2f2", nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `http://localhost:8080/3g2gf2f2`,
		},
		{
			name:      "conflict",
			inputBody: strings.NewReader(`http://yandex.ru`),
			inputURL:  "http://yandex.ru",
			mockBehavior: func(r *mock_service.MockShortener, url string) {
				r.EXPECT().AddShortURL(gomock.Any(), url).Return("", &pq.Error{Code: models.PqDuplicateErr})
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

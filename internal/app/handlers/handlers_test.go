package handlers

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/lib/pq"

	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/deleter"
	"github.com/koteyye/shortener/internal/app/service"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/koteyye/shortener/internal/app/models"
	mockservice "github.com/koteyye/shortener/internal/app/storage/mocks"
)

const (
	testSecretKey = "super_secret_key"
	testURL       = "http://yandex.ru"

	baseURL       = "http://localhost:8080"
	batch         = "/batch"
	urlListByUser = "/api/user/urls"
	jsonShortener = "/api/shorten"
	ping          = "/ping"
	stats         = "/api/internal/stats"
)

func TestHandlers_NewHandlers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testService := service.Service{}
		handler := NewHandlers(&testService, &zap.SugaredLogger{}, testSecretKey, &deleter.Deleter{}, &net.IPNet{})

		assert.Equal(t, &Handlers{
			services:  &testService,
			logger:    &zap.SugaredLogger{},
			secretKey: testSecretKey,
			worker:    &deleter.Deleter{},
			subnet:    &net.IPNet{},
		}, handler)
	})
}

// testInitHandler инициализация тестового обработчика
func testInitHandler(t *testing.T) (*Handlers, *mockservice.MockURLStorage) {
	c := gomock.NewController(t)
	defer c.Finish()

	repo := mockservice.NewMockURLStorage(c)
	service := service.NewService(repo, &config.Shortener{Listen: "http://localhost:8081"}, &zap.SugaredLogger{})

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	log := *logger.Sugar()

	handler := NewHandlers(service, &log, testSecretKey, &deleter.Deleter{}, &net.IPNet{})

	return handler, repo
}

func generateUnitKey() string {
	t := time.Now().UnixNano()
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(t)))
}

func TestHandlers_Batch(t *testing.T) {
	testRequest := `[
		{
			"correlation_id": "afd90f2c-b0df-4873-8ded-62d8e99593ba",
			"original_url": "http://lcpjtoddpyyp.yandex/sjkhh"
		}
		]`

	t.Run("batch", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			h, s := testInitHandler(t)

			r := httptest.NewRequest(http.MethodPost, batch, strings.NewReader(testRequest))
			r.Header.Add("Content-Type", ctApplicationJSON)
			w := httptest.NewRecorder()

			var testOriginURLList []*models.URLList
			err := json.Unmarshal([]byte(testRequest), &testOriginURLList)
			assert.NoError(t, err)

			userID := uuid.NewString()
			ctx := context.WithValue(r.Context(), userIDKey, userID)

			s.EXPECT().AddURL(gomock.Any(), gomock.Any(), "http://lcpjtoddpyyp.yandex/sjkhh", userID).Return(error(nil))
			s.EXPECT().GetShortURL(gomock.Any(), "http://lcpjtoddpyyp.yandex/sjkhh").Return(generateUnitKey(), error(nil))

			h.Batch(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusCreated, w.Code)
		})
		t.Run("conflict", func(t *testing.T) {
			h, s := testInitHandler(t)

			r := httptest.NewRequest(http.MethodPost, "/batch", strings.NewReader(testRequest))
			r.Header.Add("Content-Type", ctApplicationJSON)
			w := httptest.NewRecorder()

			var testOriginURLList []*models.URLList
			err := json.Unmarshal([]byte(testRequest), &testOriginURLList)
			assert.NoError(t, err)

			userID := uuid.NewString()
			ctx := context.WithValue(r.Context(), userIDKey, userID)

			testPQErr := &pq.Error{Code: models.PqDuplicateErr}

			s.EXPECT().AddURL(gomock.Any(), gomock.Any(), "http://lcpjtoddpyyp.yandex/sjkhh", userID).Return(error(testPQErr))
			s.EXPECT().GetShortURL(gomock.Any(), "http://lcpjtoddpyyp.yandex/sjkhh").Return(generateUnitKey(), error(nil))

			h.Batch(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusConflict, w.Code)
			assert.Equal(t, ctApplicationJSON, w.Header().Get("Content-Type"))
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
			r.Header.Add("Content-Type", ctApplicationJSON)
			w := httptest.NewRecorder()

			var testOriginURLList []*models.URLList
			err := json.Unmarshal([]byte(testRequest), &testOriginURLList)
			assert.NoError(t, err)

			userID := uuid.NewString()
			ctx := context.WithValue(r.Context(), userIDKey, userID)

			h.Batch(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.Equal(t, ctApplicationJSON, w.Header().Get("Content-Type"))
		})
	})
}

func TestHandlers_GetURLsByUser(t *testing.T) {
	t.Run("get_urls_by_user", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			h, s := testInitHandler(t)
			r := httptest.NewRequest(http.MethodPost, urlListByUser, nil)
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())

			repoRes := make([]*models.URLList, 1)
			repoRes[0] = &models.URLList{
				URL:      "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv",
				ShortURL: "j90iejf9032jf923jf0923fj3029f==",
			}
			s.EXPECT().GetURLByUser(gomock.Any(), userID.String()).Return(repoRes, error(nil))

			h.GetURLsByUser(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, ctApplicationJSON, w.Header().Get("Content-Type"))
		})
		t.Run("no content", func(t *testing.T) {
			h, s := testInitHandler(t)
			r := httptest.NewRequest(http.MethodPost, urlListByUser, nil)
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())

			s.EXPECT().GetURLByUser(gomock.Any(), userID.String()).Return(nil, error(sql.ErrNoRows))

			h.GetURLsByUser(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusNoContent, w.Code)
			assert.Equal(t, ctApplicationJSON, w.Header().Get("Content-Type"))
		})
	})
}

func TestHandlers_JSONShortenURL(t *testing.T) {
	testRequest := `{"url": "http://yandex.ru"}`
	t.Run("JSONShortenURL", func(t *testing.T) {
		t.Run("succes", func(t *testing.T) {
			h, s := testInitHandler(t)
			r := httptest.NewRequest(http.MethodPost, jsonShortener, strings.NewReader(testRequest))
			r.Header.Add("Content-Type", ctApplicationJSON)
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())

			s.EXPECT().AddURL(gomock.Any(), gomock.Any(), testURL, userID.String()).Return(error(nil))

			h.JSONShortenURL(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusCreated, w.Code)
			assert.Equal(t, ctApplicationJSON, w.Header().Get("Content-Type"))
		})
		t.Run("conflict", func(t *testing.T) {
			h, s := testInitHandler(t)
			r := httptest.NewRequest(http.MethodPost, jsonShortener, strings.NewReader(testRequest))
			r.Header.Add("Content-Type", ctApplicationJSON)
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())

			testPQErr := &pq.Error{Code: models.PqDuplicateErr}
			s.EXPECT().AddURL(gomock.Any(), gomock.Any(), testURL, userID.String()).Return(error(testPQErr))
			s.EXPECT().GetShortURL(gomock.Any(), testURL).Return(generateUnitKey(), error(nil))

			h.JSONShortenURL(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusConflict, w.Code)
			assert.Equal(t, ctApplicationJSON, w.Header().Get("Content-Type"))
		})
	})
}

func TestHandlers_ShortenURL(t *testing.T) {
	t.Run("shortenURL", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			h, s := testInitHandler(t)
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testURL))
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())
			s.EXPECT().AddURL(gomock.Any(), gomock.Any(), testURL, userID.String()).Return(error(nil))

			h.ShortenURL(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusCreated, w.Code)
			assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
		})
	})
}

func TestHandlers_Ping(t *testing.T) {
	t.Run("Ping", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			h, s := testInitHandler(t)
			r := httptest.NewRequest(http.MethodGet, ping, nil)
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())
			s.EXPECT().GetDBPing(gomock.Any()).Return(error(nil))

			h.Ping(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("error", func(t *testing.T) {
			h, s := testInitHandler(t)
			r := httptest.NewRequest(http.MethodGet, ping, nil)
			w := httptest.NewRecorder()

			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())
			s.EXPECT().GetDBPing(gomock.Any()).Return(errors.New("err"))

			h.Ping(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusInternalServerError, w.Code)
		})
	})
}

func TestHandlers_GetOriginalURL(t *testing.T) {
	t.Run("get_original_URL", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			shortURL := generateUnitKey()
			h, s := testInitHandler(t)
			r := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)

			w := httptest.NewRecorder()
			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())
			s.EXPECT().GetURL(gomock.Any(), gomock.Any()).Return(&models.SingleURL{URL: testURL}, error(nil))

			h.GetOriginalURL(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
		})
		t.Run("gone", func(t *testing.T) {
			shortURL := generateUnitKey()
			h, s := testInitHandler(t)
			r := httptest.NewRequest(http.MethodGet, "/"+shortURL, nil)

			w := httptest.NewRecorder()
			userID := uuid.New()
			ctx := context.WithValue(r.Context(), userIDKey, userID.String())
			s.EXPECT().GetURL(gomock.Any(), gomock.Any()).Return(&models.SingleURL{IsDeleted: true}, error(nil))

			h.GetOriginalURL(w, r.WithContext(ctx))

			assert.Equal(t, http.StatusGone, w.Code)
		})
	})
}

func TestHandlers_GetStats(t *testing.T) {
	t.Run("get stats", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockservice.NewMockURLStorage(c)
			service := service.NewService(repo, &config.Shortener{Listen: "http://localhost:8081"}, &zap.SugaredLogger{})

			logger, err := zap.NewDevelopment()
			if err != nil {
				panic(err)
			}
			defer logger.Sync()
			log := *logger.Sugar()
			cfg := config.Config{TrustSubnet: "127.0.0.1/24"}
			subnet, err := cfg.CIDR()
			assert.NoError(t, err)

			handler := NewHandlers(service, &log, testSecretKey, &deleter.Deleter{}, subnet)

			r := httptest.NewRequest(http.MethodGet, stats, nil)
			w := httptest.NewRecorder()

			repo.EXPECT().GetCount(gomock.Any()).Return(50, 20, error(nil))

			handler.GetStats(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
		})
		t.Run("failed", func(t *testing.T) {
			h, s := testInitHandler(t)

			r := httptest.NewRequest(http.MethodGet, stats, nil)
			w := httptest.NewRecorder()

			s.EXPECT().GetCount(gomock.Any()).Return(0, 0, errors.New("can't get stats from storage"))

			h.GetStats(w, r)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	})
}

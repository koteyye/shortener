package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/service"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Тестовый конфиг
var cfg = config.Config{
	Server: &config.Server{
		BaseURL: "/",
		Listen:  "localhost:8080",
	},
	Shortener: &config.Shortener{
		Listen: "http://localhost:8080",
	},
}

func TestHandlers_ShortenerURL(t *testing.T) {
	testDir := t.TempDir()
	file, _ := os.CreateTemp(testDir, "db")

	type want struct {
		statusCodePOST int
		statusCodeGET  int
		locationHeader string
		wantErr        error
	}
	tests := []struct {
		name        string
		methodType  string
		request     string
		requestBody io.Reader
		want        want
	}{
		{
			name:        "success",
			request:     "/",
			requestBody: strings.NewReader("https://practicum.yandex.ru/"),
			want: want{
				statusCodePOST: 201,
				statusCodeGET:  307,
				locationHeader: "https://practicum.yandex.ru/",
			},
		},
		{
			name:        "null body",
			request:     "/",
			requestBody: strings.NewReader(""),
			want: want{
				statusCodePOST: 400,
				wantErr:        service.ErrNullRequestBody,
			},
		},
		{
			name:        "invalid URL",
			request:     "/",
			requestBody: strings.NewReader("practicum.yandex.ru/"),
			want: want{
				statusCodePOST: 400,
				wantErr:        service.ErrInvalidRequestBodyURL,
			},
		},
	}

	storages := storage.NewURLHandle(file.Name())
	services := service.NewService(storages, cfg.Shortener)
	h := NewHandlers(services, zap.SugaredLogger{})
	hostName := cfg.Shortener.Listen + "/"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Тест POST
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, test.request, test.requestBody)
			h.ShortenerURL(c)

			result := w.Result()
			defer result.Body.Close()
			//Преобразуем тело для проверки
			body, _ := io.ReadAll(result.Body)

			if result.StatusCode == 201 {
				regCheck := strings.Contains(string(body), hostName)
				assert.Equal(t, true, regCheck)
				shortURL := strings.TrimPrefix(string(body), hostName)
				assert.Equal(t, test.want.statusCodePOST, result.StatusCode)
				wGET := httptest.NewRecorder()
				c2, _ := gin.CreateTestContext(wGET)
				c2.Request = httptest.NewRequest(http.MethodGet, test.request, nil)
				c2.AddParam("id", shortURL)
				h.LongerURL(c2)
				resultGET := wGET.Result()
				defer resultGET.Body.Close()
				assert.Equal(t, test.want.statusCodeGET, resultGET.StatusCode)
				assert.Equal(t, test.want.locationHeader, resultGET.Header.Get("Location"))
			} else {
				assert.Equal(t, test.want.statusCodePOST, result.StatusCode)
				assert.EqualError(t, test.want.wantErr, test.want.wantErr.Error(), result.Body)
			}
		})
	}
}

func TestHandlers_LongerURL(t *testing.T) {
	testDir := t.TempDir()
	file, _ := os.CreateTemp(testDir, "db")
	type want struct {
		statusCode     int
		locationHeader string
		wantErr        error
	}
	tests := []struct {
		name   string
		params string
		want
	}{
		{
			name:   "success",
			params: "MTY5NDAzNTIwNjI4NjQyNzIwOQ==",
			want: want{
				statusCode:     307,
				locationHeader: "https://practicum.yandex.ru/",
			},
		},
		{
			name:   "invalid param",
			params: "MTY5NDAzNTIwNjI4NjQyNzIds==",
			want: want{
				statusCode: 400,
				wantErr:    storage.ErrNotFound,
			},
		},
		{
			name: "not params",
			want: want{
				statusCode: 400,
			},
		},
	}
	storages := storage.NewURLHandle(file.Name())
	services := service.NewService(storages, cfg.Shortener)
	h := NewHandlers(services, zap.SugaredLogger{})

	//Предварительно добавляется валидное значение в Storage
	storages.AddURL("MTY5NDAzNTIwNjI4NjQyNzIwOQ==", "https://practicum.yandex.ru/")
	for _, test := range tests {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/"+test.params, nil)
		c.AddParam("id", test.params)
		h.LongerURL(c)
		result := w.Result()
		defer result.Body.Close()
		assert.Equal(t, test.statusCode, result.StatusCode)
		if test.statusCode == 200 {
			assert.Equal(t, test.locationHeader, result.Header.Get("Location"))
		}
		if test.wantErr != nil {
			assert.EqualError(t, test.want.wantErr, test.want.wantErr.Error(), result.Body)
		}
	}
}

func TestHandlers_ShortenerURLJSON(t *testing.T) {
	testDir := t.TempDir()
	file, err := os.CreateTemp(testDir, "db")
	require.NoError(t, err)
	type want struct {
		statusCodePOST int
		statusCodeGET  int
		locationHeader string
		contentType    string
		wantErr        error
	}
	tests := []struct {
		name        string
		methodType  string
		request     string
		requestBody models.LongURL
		want        want
	}{
		{
			name:        "success",
			request:     "/api/shorten",
			requestBody: models.LongURL{URL: "https://practicum.yandex.ru/"},
			want: want{
				statusCodePOST: 201,
				statusCodeGET:  307,
				locationHeader: "https://practicum.yandex.ru/",
				contentType:    "application/json",
			},
		},
		{
			name:        "null body",
			request:     "/",
			requestBody: models.LongURL{URL: ""},
			want: want{
				statusCodePOST: 400,
				wantErr:        service.ErrNullRequestBody,
				contentType:    "application/json",
			},
		},
		{
			name:        "invalid URL",
			request:     "/",
			requestBody: models.LongURL{URL: "practicum.yandex.ru/"},
			want: want{
				statusCodePOST: 400,
				wantErr:        service.ErrInvalidRequestBodyURL,
				contentType:    "application/json",
			},
		},
	}

	storages := storage.NewURLHandle(file.Name())
	services := service.NewService(storages, cfg.Shortener)
	h := NewHandlers(services, zap.SugaredLogger{})
	hostName := cfg.Shortener.Listen + "/"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//Тест POST
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			req, _ := json.Marshal(test.requestBody)
			reader := bytes.NewReader(req)
			c.Request = httptest.NewRequest(http.MethodPost, test.request, reader)
			h.ShortenerURLJSON(c)

			result := w.Result()
			defer result.Body.Close()
			//Преобразуем тело для проверки
			body, _ := io.ReadAll(result.Body)

			if result.StatusCode == 201 {
				var bodyShortURL models.ShortURL
				_ = json.Unmarshal(body, &bodyShortURL)
				regCheck := strings.Contains(bodyShortURL.Result, hostName)
				assert.Equal(t, true, regCheck)
				shortURL := strings.TrimPrefix(bodyShortURL.Result, hostName)
				assert.Equal(t, test.want.statusCodePOST, result.StatusCode)
				assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
				wGET := httptest.NewRecorder()
				c2, _ := gin.CreateTestContext(wGET)
				c2.Request = httptest.NewRequest(http.MethodGet, test.request, nil)
				c2.AddParam("id", shortURL)
				h.LongerURL(c2)
				resultGET := wGET.Result()
				defer resultGET.Body.Close()
				assert.Equal(t, test.want.statusCodeGET, resultGET.StatusCode)
				assert.Equal(t, test.want.locationHeader, resultGET.Header.Get("Location"))
			} else {
				assert.Equal(t, test.want.statusCodePOST, result.StatusCode)
				assert.EqualError(t, test.want.wantErr, test.want.wantErr.Error(), result.Body)
			}
		})
	}
}

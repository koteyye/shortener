package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/service"
	"github.com/koteyye/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
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
		BaseURL: "/",
		Listen:  "http://localhost:8080",
	},
}

func TestHandlers_ShortenerURL(t *testing.T) {
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

	storages := storage.NewURLHandle()
	services := service.NewService(storages, cfg.Shortener)
	h := NewHandlers(services)
	hostName := cfg.Shortener.Listen + cfg.Shortener.BaseURL
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
				regCheck, _ := regexp.Match(hostName, body)
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
	storages := storage.NewURLHandle()
	services := service.NewService(storages, cfg.Shortener)
	h := NewHandlers(services)

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

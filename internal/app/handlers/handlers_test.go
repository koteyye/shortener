package handlers

import (
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

const HostName = "http://localhost:8080/"

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
			requestBody: strings.NewReader("ofjewpogjkewpo"),
			want: want{
				statusCodePOST: 400,
				wantErr:        service.ErrInvalidRequestBodyURL,
			},
		},
	}
	storages := storage.NewURLHandle()
	services := service.NewService(storages)
	h := NewHandlers(services)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var shortURL string

			//Тест POST
			request := httptest.NewRequest(http.MethodPost, test.request, test.requestBody)
			w := httptest.NewRecorder()
			h.ShortenerURL(w, request)

			result := w.Result()

			//Преобразуем тело для проверки
			body, _ := io.ReadAll(result.Body)

			if result.StatusCode == 201 {
				regCheck, _ := regexp.Match(`(http://localhost:8080/)`, body)
				assert.Equal(t, true, regCheck)
				shortURL = strings.TrimLeft(string(body), HostName)
				assert.Equal(t, test.want.statusCodePOST, result.StatusCode)
				requestGET := httptest.NewRequest(http.MethodGet, test.request+shortURL, nil)
				wGET := httptest.NewRecorder()
				h.ShortenerURL(wGET, requestGET)
				resultGET := wGET.Result()
				assert.Equal(t, test.want.statusCodeGET, resultGET.StatusCode)
				assert.Equal(t, test.want.locationHeader, resultGET.Header.Get("Location"))
			} else {
				assert.Equal(t, test.want.statusCodePOST, result.StatusCode)
				assert.Error(t, test.want.wantErr, body)
			}
		})
	}
}

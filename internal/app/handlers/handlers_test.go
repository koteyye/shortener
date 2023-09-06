package handlers

import (
	"github.com/gin-gonic/gin"
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
				regCheck, _ := regexp.Match(`(http://localhost:8080/)`, body)
				assert.Equal(t, true, regCheck)
				shortURL := strings.TrimPrefix(string(body), HostName)
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

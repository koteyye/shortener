package handlers

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/koteyye/shortener/config"
	"github.com/koteyye/shortener/internal/app/models"
	"github.com/koteyye/shortener/internal/app/service"
	mock_service "github.com/koteyye/shortener/internal/app/service/mocks"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func TestHandlers_Batch(t *testing.T) {
	type mockBehavior func(r *mock_service.MockShortener, list []*models.OriginURLList)
	tests := []struct {
		name                 string
		inputBody            io.Reader
		inputOriginURLList   []*models.OriginURLList
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "created",
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
					ID:        "5993975-38fc-4f95-8234-0639079194cf",
					OriginURL: "http://sd37z.ru/klrotsvqdpjaj/hs0jfw6xiiv",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			r := httptest.NewRequest(http.MethodPost, "/batch", test.inputBody)
			w := httptest.NewRecorder()
			ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
			defer cancel()

			repo := mock_service.NewMockShortener(c)

			gomock.InOrder(repo.EXPECT().Batch(ctx, test.inputOriginURLList).Return(nil, nil))

			services := service.Service{Shortener: repo}
			handler := Handlers{services: &services}

			handler.Batch(w, r)

			result := w.Result()
			defer result.Body.Close()

			assert.Equal(t, w.Code, test.expectedStatusCode)
		})
	}
}

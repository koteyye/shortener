package handlers

import (
	"bytes"
	"github.com/koteyye/shortener/internal/app/service"
	"net/http"
	"strings"
)

type Handlers struct {
	services *service.Service
}

func NewHandlers(services *service.Service) *Handlers {
	return &Handlers{services: services}
}

func (h Handlers) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.shortenerURL)
	return mux
}

func (h Handlers) shortenerURL(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodPost:
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			newResponse(res, http.StatusBadRequest, "не удалось получить значение запроса")
		}
		respBytes := buf.String()

		result, err := h.services.ShortURL(respBytes)
		if err != nil {
			newResponse(res, http.StatusBadRequest, err.Error())
		}
		newResponse(res, http.StatusCreated, result)
	case http.MethodGet:
		url := req.URL.String()
		if id := strings.TrimLeft(url, "/"); id != "" {
			resURL, err := h.services.LongURL(id)
			if err != nil {
				newResponse(res, http.StatusBadRequest, "нет такой ссылки")
			}
			newResponse(res, http.StatusTemporaryRedirect, "Location: "+resURL)
		} else {
			newResponse(res, http.StatusBadRequest, "в запросе нет сокращенной ссылки")
		}

	default:
		newResponse(res, http.StatusBadRequest, "не тот метод")
	}
}

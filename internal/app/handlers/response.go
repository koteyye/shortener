package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/koteyye/shortener/internal/app/models"
)

const (
	ctApplicationJSON = "application/json"
)

type errorJSON struct {
	Message string `json:"Message"`
}

type resultJSON struct {
	Result string `json:"Result"`
}

func mapToStringResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(message))
}

func mapShortURLToJSONResponse(w http.ResponseWriter, statusCode int, result string) {
	rawResponse, err := json.Marshal(&resultJSON{Result: result})
	if err != nil {
		mapToStringResponse(w, http.StatusBadRequest, err.Error())
	}
	w.Header().Add("Content-type", ctApplicationJSON)
	w.WriteHeader(statusCode)
	w.Write(rawResponse)
}

func mapErrorToResponse(w http.ResponseWriter, r *http.Request, statusCode int, msg string) {
	contentType := r.Header.Get("Content-Type")
	w.WriteHeader(statusCode)
	if contentType == ctApplicationJSON {
		rawResponse, err := json.Marshal(&errorJSON{Message: msg})
		if err != nil {
			mapToStringResponse(w, http.StatusBadRequest, err.Error())
		}
		w.Header().Add("Content-type", ctApplicationJSON)
		w.Write(rawResponse)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(msg))
}

func mapErrorToJSONResponse(w http.ResponseWriter, statusCode int, msg string) {
	rawResponse, err := json.Marshal(&errorJSON{Message: msg})
	if err != nil {
		mapToStringResponse(w, http.StatusBadRequest, err.Error())
	}
	w.Header().Add("Content-type", ctApplicationJSON)
	w.WriteHeader(statusCode)
	w.Write(rawResponse)
}

func mapURLListToJSONResponse(w http.ResponseWriter, statusCode int, result []*models.URLList) {
	rawResponse, err := json.Marshal(result)
	if err != nil {
		mapToStringResponse(w, http.StatusBadRequest, err.Error())
	}
	w.Header().Add("Content-type", ctApplicationJSON)
	w.WriteHeader(statusCode)
	w.Write(rawResponse)
}

func mapAllURLsToJSONResponse(w http.ResponseWriter, statusCode int, result []*models.URLList) {
	rawResponse, err := json.Marshal(result)
	if err != nil {
		mapToStringResponse(w, http.StatusBadRequest, err.Error())
	}
	w.Header().Add("Content-type", ctApplicationJSON)
	w.WriteHeader(statusCode)
	w.Write(rawResponse)
}

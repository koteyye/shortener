package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LoggingResponseWriter структура логирования запроса.
type LoggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

type log struct {
	URI        string `json:"uri"`
	Method     string `json:"method"`
	Duration   int64  `json:"duration"`
	StatusCode int    `json:"statusCode"`
	Size       int    `json:"size"`
	ErrMsg     string `json:"errMsg,omitempty"`
}

type responseData struct {
	status int
	size   int
}

// Write записать байты в http-ответ.
func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader записать заголовок http-ответа.
func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func marshalJSON(s *log) ([]byte, error) {
	m, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("не удалось сериализовать лог: %v", err)
	}
	return m, nil
}

// Logging логирование http-запроса.
func (h Handlers) Logging(next http.Handler) http.Handler {
	logFN := func(res http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := LoggingResponseWriter{
			ResponseWriter: res,
			responseData:   responseData,
		}
		next.ServeHTTP(&lw, r)

		duration := time.Since(start).Milliseconds()
		logItem, err := marshalJSON(&log{
			URI:        r.RequestURI,
			Method:     r.Method,
			Duration:   duration,
			StatusCode: responseData.status,
			Size:       responseData.size,
		})
		if err != nil {
			mapErrorToResponse(res, r, http.StatusInternalServerError, err.Error())
			return
		}
		h.logger.Infow("HTTP Request", "event", string(logItem))
	}
	return http.HandlerFunc(logFN)
}

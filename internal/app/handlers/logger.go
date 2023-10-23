package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

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

func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func marshalJSON(s *log) []byte {
	m, err := json.Marshal(s)
	if err != nil {

	}

	return m
}

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

		duration := time.Since(start).Nanoseconds()

		h.logger.Infow("HTTP Request", "event", string(marshalJSON(&log{
			URI:        r.RequestURI,
			Method:     r.Method,
			Duration:   duration,
			StatusCode: responseData.status,
			Size:       responseData.size,
		})))
	}
	return http.HandlerFunc(logFN)
}

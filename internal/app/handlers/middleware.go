package handlers

import (
	"compress/gzip"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoggingResponseWriter struct {
	gin.ResponseWriter
	responseData *responseData
}

type compressWriter struct {
	w  gin.ResponseWriter
	zw *gzip.Writer
}

type log struct {
	URI        string      `json:"uri"`
	Method     string      `json:"method"`
	Duration   int64       `json:"duration"`
	StatusCode int         `json:"statusCode"`
	Size       int         `json:"size"`
	Headers    http.Header `json:"headers"`
}

type responseData struct {
	status int
	size   int
}

func (h Handlers) Compressing() gin.HandlerFunc {
	return Compress()
}

func (h Handlers) WithLogging() gin.HandlerFunc {
	logFn := func(c *gin.Context) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := &LoggingResponseWriter{
			ResponseWriter: c.Writer,
			responseData:   responseData,
		}
		c.Writer = lw
		c.Next()
		duration := time.Since(start).Nanoseconds()

		h.logger.Infow("HTTP Request", "event", string(marshalJSON(&log{
			URI:        c.Request.RequestURI,
			Method:     c.Request.Method,
			Duration:   duration,
			StatusCode: c.Writer.Status(),
			Size:       c.Writer.Size(),
			Headers:    c.Request.Header,
		})))
	}
	return logFn
}

func marshalJSON(s *log) []byte {
	m, err := json.Marshal(s)
	if err != nil {
		panic(err.Error())
	}

	return m
}

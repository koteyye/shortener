package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

type log struct {
	Uri        string        `json:"uri"`
	Method     string        `json:"method"`
	StatusCode int           `json:"statusCode"`
	Duration   time.Duration `json:"duration"`
	Size       int           `json:"size"`
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	gin.ResponseWriter
	responseData *responseData
}

func (h Handlers) WithLogging() gin.HandlerFunc {
	logFn := func(c *gin.Context) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lw := &loggingResponseWriter{
			ResponseWriter: c.Writer,
			responseData:   responseData,
		}
		c.Writer = lw
		c.Next()
		duration := time.Since(start)

		h.logger.Infow("HTTP Request", "event", string(marshalJSON(&log{
			Uri:        c.Request.RequestURI,
			Method:     c.Request.Method,
			Duration:   duration,
			StatusCode: c.Writer.Status(),
			Size:       c.Writer.Size(),
		})))
	}

	//h.logger.Infoln(
	//	"uri", c.Request.RequestURI,
	//	"method", c.Request.Method,
	//	"duration", duration,
	//	"statusCode", c.Writer.Status(),
	//	"size", c.Writer.Size(),
	//)
	return logFn
}

func marshalJSON(s *log) []byte {
	m, err := json.Marshal(s)
	if err != nil {
		panic(err.Error())
	}

	return m
}

package handlers

import (
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/http"
)

func Compress() gin.HandlerFunc {
	return func(c *gin.Context) {

		//
		contentEncoding := c.Request.Header.Get("Content-Encoding")
		if contentEncoding == "gzip" {
			gz, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				newJSONResponse(c, http.StatusBadRequest, err)
			}
			defer gz.Close()

			c.Request.Body = gz
		}

		acceptGzip := c.Request.Header.Values("Accept-Encoding")
		acceptContent := []string{"application/json", "text/html"}
		isAcceptGzip := slices.Contains(acceptGzip, "gzip") && slices.Contains(acceptContent, c.Request.Header.Get("Content-Type"))
		if isAcceptGzip {
			gw := gzip.NewWriter(c.Writer)
			gw.Reset(c.Writer)
			c.Writer = &gzipWriter{
				ResponseWriter: c.Writer,
				writer:         gw,
			}
			c.Header("Content-Encoding", "gzip")
			defer c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
		}

		c.Next()

	}
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	g.Header().Del("Content-Length")
	return g.writer.Write(data)
}

func (g *gzipWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

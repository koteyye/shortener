package handlers

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/http"
)

func Compress() gin.HandlerFunc {
	return func(c *gin.Context) {

		contentEncoding := c.Request.Header.Get("Content-Encoding")
		if contentEncoding == "gzip" {
			gr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				newJSONResponse(c, http.StatusBadRequest, err)
				return
			}
			defer gr.Close()

			c.Request.Body = gr
		}

		acceptEncoding := c.Request.Header.Values("Accept-Encoding")
		supportGzip := slices.Contains(acceptEncoding, "gzip")
		if supportGzip {
			gw := gzip.NewWriter(c.Writer)
			c.Writer = &gzipWriter{
				ResponseWriter: c.Writer,
				writer:         gw,
			}
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

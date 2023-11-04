package handlers

import (
	"compress/gzip"
	"fmt"
	"golang.org/x/exp/slices"
	"net/http"
)

func (h Handlers) Compress(next http.Handler) http.Handler {
	compressFn := func(res http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		if contentEncoding == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				mapErrorToResponse(res, r, http.StatusBadRequest, fmt.Sprintf("gzip newreader: %v", err))
				return
			}
			defer gz.Close()

			r.Body = gz
		}

		acceptGzip := r.Header.Values("Accept-Encoding")
		acceptContent := []string{"application/json", "text/html"}
		isAcceptGzip := slices.Contains(acceptGzip, "gzip") && slices.Contains(acceptContent, r.Header.Get("Content-Type"))
		if isAcceptGzip {
			gw := gzip.NewWriter(res)
			gw.Reset(res)
			res = &gzipWriter{
				ResponseWriter: res,
				writer:         gw,
			}
			res.Header().Add("Content-Encoding", "gzip")
			defer func() {
				gw.Close()
				res.Header().Add("Content-Length", fmt.Sprint(res.Header().Get("size")))
			}()
		}
		next.ServeHTTP(res, r)
	}
	return http.HandlerFunc(compressFn)
}

type gzipWriter struct {
	http.ResponseWriter
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

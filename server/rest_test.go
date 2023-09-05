package server

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func mainPage(res http.ResponseWriter, req *http.Request) {
	body := fmt.Sprintf("Method: %s\r\n", req.Method)
	body += "Header ===============\r\n"
	for k, v := range req.Header {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	body += "Query parameters ===============\r\n"
	if err := req.ParseForm(); err != nil {
		res.Write([]byte(err.Error()))
		return
	}
	for k, v := range req.Form {
		body += fmt.Sprintf("%s: %v\r\n", k, v)
	}
	res.Write([]byte(body))
}

func TestServer_Run(t *testing.T) {
	s := Server{}

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage)

	go func() {
		_ = s.Run("8081", mux)
	}()
	defer s.Shutdown(context.Background())
	r, _ := http.NewRequest(http.MethodGet, "http://localhost:8081", nil)
	res, err := http.DefaultClient.Do(r)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
	}

}

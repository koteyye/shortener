package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kabukky/httpscerts"
	"go.uber.org/zap"
)

const (
	certFile = "cert.pem"
	keyFile  = "key.pem"
)

// Server определяет структуру сервера.
type Server struct {
	httpServer *http.Server
}

// Run запускает сервер.
func (s *Server) Run(enableHTTPS bool, host string, handler http.Handler) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := *logger.Sugar()

	s.httpServer = &http.Server{
		Addr:    host,
		Handler: handler,
	}

	sugar.Info("starting server")
	if enableHTTPS {
		err := httpscerts.Check(certFile, keyFile)
		if err != nil {
			err = httpscerts.Generate(certFile, keyFile, host)
			if err != nil {
				return fmt.Errorf("can't generate https cert: %w", err)
			}
		}
		return s.httpServer.ListenAndServeTLS(certFile, keyFile)
	}
	return s.httpServer.ListenAndServe()
}

// Shutdown отключает сервер.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

package server

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

// Server определяет структуру сервера.
type Server struct {
	httpServer *http.Server
}

// Run запускает сервер.
func (s *Server) Run(host string, handler http.Handler) error {
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
	return s.httpServer.ListenAndServe()
}

// Shutdown отключает сервер.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

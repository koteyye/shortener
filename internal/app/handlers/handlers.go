package handlers

import (
	"github.com/go-chi/chi"
	"github.com/koteyye/shortener/internal/app/service"
	"go.uber.org/zap"
)

type Handlers struct {
	services  *service.Service
	logger    zap.SugaredLogger
	secretKey string
}

func NewHandlers(services *service.Service, logger zap.SugaredLogger, secretKey string) *Handlers {
	return &Handlers{services: services, logger: logger, secretKey: secretKey}
}

func (h Handlers) InitRoutes(baseURL string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(h.Logging)
	r.Use(h.Compress)
	r.Use(h.Authorization)
	r.Post(baseURL, h.ShortenURL)
	r.Route(baseURL, func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Use(h.Authorization)
			r.Use(h.mapParamsGetOriginalURL)
			r.Get("/", h.GetOriginalURL)
		})
		r.Get("/:id", h.GetOriginalURL)
		r.Get("/ping", h.Ping)
		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", h.JSONShortenURL)
				r.Post("/batch", h.Batch)
			})
			r.Route("/user/urls", func(r chi.Router) {
				r.Get("/", h.GetURLsByUser)
			})
		})
	})

	return r
}

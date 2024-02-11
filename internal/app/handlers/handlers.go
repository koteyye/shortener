package handlers

import (
	"net"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"

	"github.com/koteyye/shortener/internal/app/deleter"
	"github.com/koteyye/shortener/internal/app/service"

	_ "github.com/koteyye/shortener/docs"
)

// Handlers http обработик сервиса
type Handlers struct {
	services  *service.Service
	logger    *zap.SugaredLogger
	delURLch chan deleter.DeleteURL
	secretKey string
	subnet    *net.IPNet
}

// NewHandlers возвращает экземпляр http обработчика
func NewHandlers(services *service.Service, logger *zap.SugaredLogger, secretKey string, delURLch chan deleter.DeleteURL, subnet *net.IPNet) *Handlers {
	return &Handlers{services: services, logger: logger, secretKey: secretKey, delURLch: delURLch, subnet: subnet}
}

// InitRoutes инициализация роутов
func (h Handlers) InitRoutes(baseURL string) *chi.Mux {
	r := chi.NewRouter()

	r.Use(h.Logging)
	r.Use(h.Compress)
	r.Use(h.Authorization)
	r.Post(baseURL, h.ShortenURL)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), //The url pointing to API definition
	))

	r.Route(baseURL, func(r chi.Router) {
		r.Route("/{id}", func(r chi.Router) {
			r.Use(h.Authorization)
			r.Get("/", h.GetOriginalURL)
		})
		r.Get("/:id", h.GetOriginalURL)
		r.Get("/ping", h.Ping)
		r.Get("/testGraceful", h.graceful)
		r.Route("/api", func(r chi.Router) {
			r.Route("/shorten", func(r chi.Router) {
				r.Post("/", h.JSONShortenURL)
				r.Post("/batch", h.Batch)
			})
			r.Route("/user/urls", func(r chi.Router) {
				r.Get("/", h.GetURLsByUser)
				r.Delete("/", h.DeleteURLsByUser)
			})
			r.Route("/internal/stats", func(r chi.Router) {
				r.Use(h.Subnet)
				r.Get("/", h.GetStats)
			})
		})
	})

	return r
}

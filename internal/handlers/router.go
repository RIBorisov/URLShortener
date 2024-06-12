package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "shortener/internal/middleware"
	"shortener/internal/service"
)

func NewRouter(svc *service.Service) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(mw.GzipMiddleware)
	router.Use(mw.Logger)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(svc))
		r.Post("/", SaveHandler(svc))
	})
	router.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", ShortenHandler(svc))
		r.Post("/batch", BatchHandler(svc))
	})

	router.Get("/ping", PingHandler(svc))
	return router
}

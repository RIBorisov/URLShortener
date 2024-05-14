package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"shortener/internal/config"
	"shortener/internal/service"
)

func NewRouter(svc *service.Service, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(svc))
		r.Post("/", SaveHandler(svc, cfg))
	})

	return router
}

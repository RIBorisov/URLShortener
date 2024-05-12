package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"shortener/internal/config"
	"shortener/internal/storage"
)

func NewRouter(db *storage.Storage, cfg *config.Config) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(db, cfg))
		r.Post("/", SaveHandler(db, cfg))
	})

	return router
}

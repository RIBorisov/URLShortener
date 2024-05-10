package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"shortener/internal/storage"

	"shortener/internal/handlers/routes"
)

func Init(db *storage.SimpleStorage) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", routes.GetHandler(db))
		r.Post("/", routes.SaveHandler(db))
	})

	return router
}

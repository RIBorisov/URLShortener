package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"shortener/internal/handlers/routes"
)

func Init() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", routes.GetURLHandler)
		r.Post("/", routes.SaveURLHandler)
	})

	return router
}

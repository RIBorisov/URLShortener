package handlers

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "shortener/internal/middleware"
	"shortener/internal/service"
)

func NewRouter(svc *service.Service, log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(mw.New(log))
	router.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(svc))
		r.Post("/", SaveHandler(svc))
	})
	router.Post("/api/shorten", ShortenHandler(svc))

	return router
}

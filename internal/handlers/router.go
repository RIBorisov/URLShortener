package handlers

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "shortener/internal/middleware"
	"shortener/internal/service"
)

func NewRouter(ctx context.Context, svc *service.Service) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(mw.GzipMiddleware)
	router.Use(mw.Logger)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(svc))
		r.Post("/", SaveHandler(svc))
	})
	router.Post("/api/shorten", ShortenHandler(svc))
	router.Get("/ping", PingHandler(ctx, svc))
	return router
}

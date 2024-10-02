// Package handlers using for routing user requests.
package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "shortener/internal/middleware"
	"shortener/internal/service"
)

// NewRouter creates router.
func NewRouter(svc *service.Service) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(mw.Auth(svc).Middleware)
	router.Use(mw.Gzip(svc.Log).Middleware)
	router.Use(mw.Log(svc.Log).Middleware)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(svc))
		r.Post("/", SaveHandler(svc))
	})
	router.Route("/api", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", ShortenHandler(svc))
			r.Post("/batch", BatchHandler(svc))
		})
		r.Route("/user", func(r chi.Router) {
			r.Use(mw.CheckAuth(svc.Log).Middleware)
			r.Get("/urls", GetURLsHandler(svc))
		})
	})
	router.Delete("/api/user/urls", DeleteURLsHandler(svc))
	router.Get("/api/internal/stats", StatsHandler(svc))
	router.Get("/ping", PingHandler(svc))
	router.Mount("/debug", middleware.Profiler())

	return router
}

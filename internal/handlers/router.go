package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"shortener/internal/models"

	mw "shortener/internal/middleware"
	"shortener/internal/service"
)

func NewRouter(svc *service.Service) *chi.Mux {
	router := chi.NewRouter()
	currentUser := &models.User{}

	router.Use(middleware.Recoverer)
	router.Use(mw.Gzip(svc.Log).Middleware)
	router.Use(mw.Log(svc.Log).Middleware)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(svc))
		r.Post("/", SaveHandler(svc, currentUser))
	})
	router.Route("/api", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", ShortenHandler(svc, currentUser))
			r.Post("/batch", BatchHandler(svc, currentUser))
		})
		r.Route("/user", func(r chi.Router) {
			r.Get("/urls", GetURLsHandler(svc, currentUser))
		})
	})

	router.Get("/ping", PingHandler(svc))
	return router
}

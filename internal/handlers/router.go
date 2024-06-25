package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	mw "shortener/internal/middleware"
	"shortener/internal/models"
	"shortener/internal/service"
)

func NewRouter(svc *service.Service) *chi.Mux {
	router := chi.NewRouter()
	user := &models.User{}

	router.Use(middleware.Recoverer)
	router.Use(mw.Auth(svc.Log, user).Middleware)
	router.Use(mw.Gzip(svc.Log).Middleware)
	router.Use(mw.Log(svc.Log).Middleware)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", GetHandler(svc))
		r.Post("/", SaveHandler(svc, user))
	})
	router.Route("/api", func(r chi.Router) {
		r.Route("/shorten", func(r chi.Router) {
			r.Post("/", ShortenHandler(svc, user))
			r.Post("/batch", BatchHandler(svc, user))
		})
		r.Route("/user", func(r chi.Router) {
			r.Use(mw.CheckAuth(svc.Log, user).Middleware)
			r.Get("/urls", GetURLsHandler(svc, user))
		})
	})
	router.Delete("/api/user/urls", DeleteURLsHandler(svc, user))
	router.Get("/ping", PingHandler(svc))
	return router
}

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"shortener/internal/config"
	"shortener/internal/handlers/routes"
)

func main() {
	cfg := config.MustLoad()
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", routes.GetURLHandler)
		r.Post("/", routes.SaveURLHandler)
	})

	log.Println("Server started on port 8080")
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

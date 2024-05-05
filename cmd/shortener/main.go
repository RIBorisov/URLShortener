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
	parseFlags()
	cfg := config.MustLoad()

	// settings server address
	if flagRunAddr != "" {
		cfg.HTTPServer.Address = flagRunAddr
	}
	router := chi.NewRouter()

	// setting middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", routes.GetURLHandler)
		r.Post("/", routes.SaveURLHandler)
	})

	//srv := &http.Server{
	//	Addr:         cfg.HTTPServer.Address,
	//	Handler:      router,
	//	ReadTimeout:  cfg.Timeout,
	//	WriteTimeout: cfg.Timeout,
	//	IdleTimeout:  cfg.IdleTimeout,
	//}
	log.Printf("\nServer running on %s\n", cfg.HTTPServer.Address)

	//if err := srv.ListenAndServe(); err != nil {
	if err := http.ListenAndServe(cfg.HTTPServer.Address, router); err != nil {
		log.Fatalf("Got unexpected error, details: %s", err)
	}
}

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	c "shortener/internal/config"
	"shortener/internal/handlers/routes"
)

func main() {
	cfg := c.LoadConfig()

	// setting routes and middlewares
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Route("/", func(r chi.Router) {
		r.Get("/{id}", routes.GetURLHandler)
		r.Post("/", routes.SaveURLHandler)
	})
	// setting server config
	//srv := &http.Server{
	//	Addr:    c.Config.ServerAddress,
	//	Handler: router,
	//}

	log.Println("server running on:", cfg.Server.ServerAddress)
	if err := http.ListenAndServe(cfg.Server.ServerAddress, router); err != nil {
		log.Fatalf("got unexpected error, details: %s", err)
	}
}

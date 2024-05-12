package main

import (
	"log"
	"net/http"
	"shortener/internal/handlers"
	"shortener/internal/storage"

	"shortener/internal/config"
	"shortener/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	// setting routes and middlewares
	db := storage.GetStorage()
	h := handlers.NewHandler(db)
	r := router.Init(h)

	// setting server config
	srv := &http.Server{
		Addr:    cfg.Server.ServerAddress,
		Handler: r,
	}

	log.Printf("server running on: %s", cfg.Server.ServerAddress)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("got unexpected error, details: %s", err)
	}
}

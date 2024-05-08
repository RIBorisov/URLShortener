package main

import (
	"log"
	"net/http"

	"shortener/internal/config"
	"shortener/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	// setting routes and middlewares
	r := router.Init()

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

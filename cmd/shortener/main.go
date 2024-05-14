package main

import (
	"log"
	"net/http"

	"shortener/internal/config"
	"shortener/internal/handlers"
	"shortener/internal/service"
	"shortener/internal/storage"
)

func main() {
	cfg := config.LoadConfig()
	db := storage.LoadStorage()
	svc := &service.Service{DB: db, BaseURL: cfg.Server.BaseURL}
	r := handlers.NewRouter(svc, cfg)

	srv := &http.Server{
		Addr:    cfg.Server.ServerAddress,
		Handler: r,
	}

	log.Printf("server running on: %s", cfg.Server.ServerAddress)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("got unexpected error, details: %s", err)
	}
}

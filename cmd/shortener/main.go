package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"shortener/internal/handlers/routes"
)

func main() {
	//cfg := config.MustLoad()

	r := chi.NewRouter()

	r.Use(
		middleware.RequestID,
		middleware.SetHeader("Content-Type", "text/plain; charset=utf-8"),
	)

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", routes.GetURLHandler)
		r.Post("/", routes.SaveURLHandler)
	})

	log.Println("Server started on port 8080")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

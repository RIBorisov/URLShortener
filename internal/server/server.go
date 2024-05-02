package server

import (
	"log"
	"net/http"
	"shortener/internal/handlers/routes"
)

func RunServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", routes.RootHandler)

	log.Println("Server started on port 8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

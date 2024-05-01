package server

import (
	h "github.com/RIBorisov/URLShortener/handlers"
	"log"
	"net/http"
)

func RunServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.RootHandler)

	log.Println("Server started on port 8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

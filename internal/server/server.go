package server

//package server
//
//import (
//	"github.com/go-chi/chi/v5"
//	"log"
//	"net/http"
//	"shortener/internal/handlers/routes"
//)
//
//func RunServer() {
//	r := chi.NewRouter()
//	r.Route("/", func(r chi.Router) {
//		r.Get("/", routes.GetURLHandler)
//		r.Post("/", routes.SaveURLHandler)
//	})
//
//	log.Println("Server started on port 8080")
//
//	err := http.ListenAndServe(":8080", r)
//	if err != nil {
//		panic(err)
//	}
//}

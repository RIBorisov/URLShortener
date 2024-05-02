package routes

import (
	"net/http"
	"shortener/internal/handlers"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetURLHandler(w, r)
	case http.MethodPost:
		SaveURLHandler(w, r)
	default:
		handlers.ErrorMethodHandler(w, []string{"POST", "GET"})
	}
}

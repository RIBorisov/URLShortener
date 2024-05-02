package handlers

import "net/http"

func RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetURLHandler(w, r)
	case http.MethodPost:
		SaveURLHandler(w, r)
	default:
		ErrorMethodHandler(w, []string{"POST", "GET"})
	}
}

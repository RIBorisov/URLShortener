package handlers

import "net/http"

func SetHeadersHandler(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

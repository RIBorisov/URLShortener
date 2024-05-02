package handlers

import (
	"net/http"
	"shortener/internal/storage"
	"strings"
)

func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortURL := strings.TrimPrefix(r.URL.Path, "/")
	SetHeadersHandler(w)

	longURL := storage.URLMap[shortURL]
	if longURL == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	originalURL := GenerateURL(r, shortURL)
	redirectToURL(w, r, longURL, originalURL)

	_, err := w.Write([]byte(longURL))
	if err != nil {
		return
	}
}

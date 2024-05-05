package routes

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"shortener/internal/handlers"
	"shortener/internal/storage"
)

func GetURLHandler(w http.ResponseWriter, r *http.Request) {
	shortLink := chi.URLParam(r, "id")
	mapper := storage.Mapper
	longLink, ok := mapper.Get(shortLink)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}

	originalURL := handlers.GenerateURL(shortLink)
	handlers.RedirectToURL(w, r, longLink, originalURL)

	_, err := w.Write([]byte(longLink))
	if err != nil {
		return
	}
}

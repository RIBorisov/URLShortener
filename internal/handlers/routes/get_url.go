package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

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

	originalURL := handlers.GetOriginalURL(shortLink)

	handlers.RedirectToURL(w, r, longLink, originalURL)

	_, err := w.Write([]byte(longLink))
	if err != nil {
		log.Printf("Error when getting ShortURL: %s", err)
		http.Error(w, "Error when saving URL", http.StatusInternalServerError)
		return
	}
}

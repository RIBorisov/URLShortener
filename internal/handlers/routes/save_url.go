package routes

import (
	"io"
	"log"
	"net/http"

	"shortener/internal/handlers"
	"shortener/internal/storage"
)

func SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	longLink, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when reading body value", http.StatusInternalServerError)
		return
	}

	mapper := storage.Mapper
	shortLink := handlers.GenerateUniqueShortLink()
	mapper.Set(shortLink, string(longLink))

	responseValue := handlers.GetOriginalURL(shortLink)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(responseValue))
	if err != nil {
		log.Printf("Error when saving URL: %s", err)
		http.Error(w, "Error when saving URL", http.StatusInternalServerError)
		return
	}
}

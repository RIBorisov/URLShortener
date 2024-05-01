package handlers

import (
	"github.com/RIBorisov/URLShortener/internal/storage"
	"io"
	"net/http"
)

func SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)

	longURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when reading body value", http.StatusBadRequest)
	}
	shortURL := "EwHXdJfB"
	storage.URLMap[shortURL] = string(longURL)

	SetHeadersHandler(w)
	responseValue := GenerateURL(r, shortURL)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(responseValue))
	if err != nil {
		return
	}
}

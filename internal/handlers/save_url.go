package handlers

import (
	"io"
	"math/rand"
	"net/http"
	"shortener/internal/storage"
)

func SaveURLHandler(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// добавить логов
			return
		}
	}(r.Body)

	longURL, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when reading body value", http.StatusBadRequest)
	}
	shortURL := generateRandomString(8)
	storage.URLMap[shortURL] = string(longURL)

	SetHeadersHandler(w)
	responseValue := GenerateURL(r, shortURL)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(responseValue))
	if err != nil {
		// добавить логов
		return
	}
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}

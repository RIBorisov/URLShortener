package routes

import (
	"io"
	"math/rand"
	"net/http"
	"shortener/internal/handlers"
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

	longLink, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error when reading body value", http.StatusBadRequest)
	}
	shortLink := generateRandomString(8)
	mapper := storage.Mapper
	mapper.Set(shortLink, string(longLink))
	handlers.SetHeadersHandler(w)
	responseValue := handlers.GenerateURL(r, shortLink)
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

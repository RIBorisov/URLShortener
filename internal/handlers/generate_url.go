package handlers

import (
	"log"
	"math/rand"
	"net/url"
	"shortener/internal/config"
	"shortener/internal/storage"
)

func GetOriginalURL(shortLink string) string {
	cfg := config.LoadConfig()
	generated, err := url.JoinPath(cfg.Server.BaseURL, shortLink)
	if err != nil {
		log.Printf("Error when generating URL %s: ", err)
	}
	return generated
}

func GenerateUniqueShortLink() string {
	var uniqString string
	cfg := config.LoadConfig()
	mapper := storage.Mapper
	// check if the string is unique
	for {
		uniqStringCandidate := generateRandomString(cfg.URL.Length)
		_, ok := mapper.Get(uniqStringCandidate)
		if !ok {
			uniqString = uniqStringCandidate
			break
		}
	}
	return uniqString
}

func generateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, length)
	// generate a random string
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}

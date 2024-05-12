package handlers

import (
	"log"
	"math/rand"
	"net/url"
	"shortener/internal/config"
)

type urlStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

func GetOriginalURL(shortLink string) string {
	cfg := config.LoadConfig()
	generated, err := url.JoinPath(cfg.Server.BaseURL, shortLink)
	if err != nil {
		log.Printf("Error when generating URL %s: ", err)
		return "" // TODO: обсудить на 1:1 (можно лучше)
	}
	return generated
}

// GenerateUniqueShortLink TODO: обсудить на 1:1 (можно лучше, не нравится передавать сюда db).
func GenerateUniqueShortLink(db urlStorage) string {
	const length = 8
	var uniqString string

	// check if the string is unique
	for {
		uniqStringCandidate := generateRandomString(length)
		_, ok := db.Get(uniqStringCandidate)
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

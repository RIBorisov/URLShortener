package service

import (
	"log"
	"math/rand"
	"net/url"
)

type urlStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

type Service struct {
	db      urlStorage
	baseURL string
}

func (h *Service) GetOriginalURL(shortLink string) string {
	// TODO(SSH)
	// cfg := config.LoadConfig() // TODO: обсудить 1:1 возможно лучше один раз инициализировать и передавать в роутер?
	// generated, err := url.JoinPath(cfg.Server.BaseURL, shortLink)
	generated, err := url.JoinPath(h.baseURL, shortLink)
	if err != nil {
		log.Printf("Error when generating URL %s: ", err)
		return "" // : обсудить на 1:1 (можно лучше)
	}
	return generated
}

// GenerateUniqueShortLink TODO: обсудить на 1:1 (можно лучше, не нравится передавать сюда db).
func (h *Service) GenerateUniqueShortLink() string {
	const length = 8
	var uniqString string

	// check if the string is unique
	for {
		uniqStringCandidate := generateRandomString(length)
		// TODO(SSH)
		_, ok := h.db.Get(uniqStringCandidate)
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

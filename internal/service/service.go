package service

import (
	"math/rand"
)

type urlStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

type Service struct {
	DB      urlStorage
	BaseURL string
}

func (s *Service) GenerateUniqueShortLink(length int) string {
	var uniqString string

	// check if the string is unique
	for {
		uniqStringCandidate := generateRandomString(length)
		_, ok := s.DB.Get(uniqStringCandidate)
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

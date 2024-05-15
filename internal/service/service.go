package service

import (
	"fmt"
	"math/rand"
)

type URLStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

type Service struct {
	DB      URLStorage
	BaseURL string
}

func (s *Service) SaveURL(long string) string {
	short := s.GenerateUniqueShortLink()
	s.DB.Save(short, long)
	return short
}

func (s *Service) GetURL(short string) (string, error) {
	long, ok := s.DB.Get(short)
	if !ok {
		return "", fmt.Errorf("not found long URL by passed short URL: %s", short)
	}
	return long, nil
}

func (s *Service) GenerateUniqueShortLink() string {
	const length = 8
	var uniqString string

	// check if the string is unique
	for {
		uniqStringCandidate := GenerateRandomString(length)
		_, ok := s.DB.Get(uniqStringCandidate)
		if !ok {
			uniqString = uniqStringCandidate
			break
		}
	}
	return uniqString
}

func GenerateRandomString(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, length)
	// generate a random string
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}

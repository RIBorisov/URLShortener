package service

import (
	"fmt"
	"math/rand"
	"shortener/internal/logger"
	"shortener/internal/storage"
)

type URLStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

type Service struct {
	DB              URLStorage
	FileStoragePath string
	BaseURL         string
}

func (s *Service) SaveURL(long string) string {
	short := s.generateUniqueShortLink()
	s.DB.Save(short, long)
	err := storage.AppendToFile(s.FileStoragePath, short, long)
	if err != nil {
		logger.Err("failed append to file %w", err)
	}
	return short
}

func (s *Service) GetURL(short string) (string, error) {
	long, ok := s.DB.Get(short)
	if !ok {
		return "", fmt.Errorf("not found long URL by passed short URL: %s", short)
	}
	return long, nil
}

func (s *Service) generateUniqueShortLink() string {
	const length = 8
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

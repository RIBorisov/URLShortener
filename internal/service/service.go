package service

import (
	"context"
	"fmt"
	"math/rand"
	"shortener/internal/models"
	"shortener/internal/storage"
)

type Service struct {
	Ctx             context.Context
	Storage         storage.URLStorage
	FileStoragePath string
	BaseURL         string
	DatabaseDSN     string
}

func (s *Service) SaveURL(long string) string {
	short := s.generateUniqueShortLink()
	s.Storage.Save(s.Ctx, short, long)
	return short
}

func (s *Service) GetURL(short string) (string, error) {
	long, ok := s.Storage.Get(s.Ctx, short)
	if !ok {
		return "", fmt.Errorf("not found long URL by passed short URL: %s", short)
	}
	return long, nil
}

func (s *Service) SaveURLs(ctx context.Context, input []models.BatchRequest) (models.BatchOut, error) {
	processed := convertData(input)
	saved, err := s.Storage.BatchSave(ctx, processed)
	if err != nil {
		return nil, fmt.Errorf("failed to batch save urls: %w", err)
	}
	return saved, nil
}

func (s *Service) generateUniqueShortLink() string {
	const length = 8
	var uniqString string

	// check if the string is unique
	for {
		uniqStringCandidate := generateRandomString(length)
		_, ok := s.Storage.Get(s.Ctx, uniqStringCandidate)
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

func convertData(input []models.BatchRequest) models.BatchIn {
	res := make(models.BatchIn, 0)
	for _, item := range input {
		res = append(res, models.BatchRequest{
			CorrelationID: item.CorrelationID,
			OriginalURL:   item.OriginalURL,
		})
	}
	return res
}

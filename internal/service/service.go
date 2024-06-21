package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/url"

	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/storage"
)

type Service struct {
	Log             *logger.Log
	Storage         storage.URLStorage
	FileStoragePath string
	BaseURL         string
	DatabaseDSN     string
}

func (s *Service) SaveURL(ctx context.Context, long string, user *models.User) (string, error) {
	short := s.generateUniqueShortLink(ctx)
	if err := s.Storage.Save(ctx, short, long, user); err != nil {
		return short, fmt.Errorf("failed save URL: %w", err)
	}
	return short, nil
}

func (s *Service) GetURL(ctx context.Context, short string) (string, error) {
	long, err := s.Storage.Get(ctx, short)
	if err != nil {
		return "", fmt.Errorf("not found long URL by passed short URL: %w", err)
	}
	return long, nil
}

func (s *Service) SaveURLs(
	ctx context.Context,
	input []models.BatchRequest,
	user *models.User,
) (models.BatchResponseArray, error) {
	processed := s.convertData(ctx, input)
	saved, err := s.Storage.BatchSave(ctx, processed, user)
	if err != nil {
		return nil, fmt.Errorf("failed to batch save urls: %w", err)
	}
	resp := make(models.BatchResponseArray, 0)
	for _, svd := range saved {
		resp = append(resp, models.BatchResponse{
			CorrelationID: svd.CorrelationID,
			ShortURL:      svd.ShortURL,
		})
	}

	return resp, nil
}

func (s *Service) generateUniqueShortLink(ctx context.Context) string {
	const length = 8
	var uniqString string

	// check if the string is unique
	for {
		uniqStringCandidate := generateRandomString(length)
		_, err := s.Storage.Get(ctx, uniqStringCandidate)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				uniqString = uniqStringCandidate
				break
			}
		}
	}
	return uniqString
}

func (s *Service) convertData(ctx context.Context, input []models.BatchRequest) models.BatchArray {
	res := make(models.BatchArray, 0)
	for _, item := range input {
		short := s.generateUniqueShortLink(ctx)
		res = append(res, models.Batch{
			CorrelationID: item.CorrelationID,
			OriginalURL:   item.OriginalURL,
			ShortURL:      short,
		})
	}
	return res
}
func (s *Service) GetUserURLs(ctx context.Context, user *models.User) (models.UserURLs, error) {
	data, err := s.Storage.GetByUserID(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed get urls by userID: %w", err)
	}
	userURLs := make(models.UserURLs, 0)
	for _, item := range data {
		short, err := url.JoinPath(s.BaseURL, "/", item.Short)
		if err != nil {
			return nil, fmt.Errorf("failed join url for short: %w", err)
		}
		userURLs = append(userURLs, models.URL{ShortURL: short, OriginalURL: item.Long})
	}
	return userURLs, nil
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

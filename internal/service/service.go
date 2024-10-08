// Package service contains the business logic of the project.
package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"shortener/internal/logger"
	"shortener/internal/models"
)

// URLStorage contains contracts for communicate with storage.
type URLStorage interface {
	Close() error
	Ping(ctx context.Context) error
	Get(ctx context.Context, shortLink string) (string, error)
	Save(ctx context.Context, shortLink, longLink string) error
	BatchSave(ctx context.Context, input models.BatchArray) (models.BatchArray, error)
	GetByUserID(ctx context.Context) ([]models.BaseRow, error)
	DeleteURLs(ctx context.Context, input models.DeleteURLs) error
	Cleanup(ctx context.Context) ([]string, error)
	ServiceStats(ctx context.Context) (models.Stats, error)
}

// Service represents the main service structure for the URL shortener.
type Service struct {
	Log             *logger.Log
	Storage         URLStorage
	FileStoragePath string
	BaseURL         string
	DatabaseDSN     string
	SecretKey       string
	TrustedSubnet   string
}

// Claims represents the claims for a JWT token.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// SaveURL saves a long URL and returns a shortened URL.
func (s *Service) SaveURL(ctx context.Context, long string) (string, error) {
	short := s.generateUniqueShortLink(ctx)
	if err := s.Storage.Save(ctx, short, long); err != nil {
		return short, fmt.Errorf("failed save URL: %w", err)
	}
	return short, nil
}

// GetURL retrieves a long URL by its short URL.
func (s *Service) GetURL(ctx context.Context, short string) (string, error) {
	long, err := s.Storage.Get(ctx, short)
	if err != nil {
		return "", fmt.Errorf("not found long URL by passed short URL: %w", err)
	}
	return long, nil
}

// SaveURLs saves multiple URLs in batch and returns the corresponding short URLs.
func (s *Service) SaveURLs(ctx context.Context, input []models.BatchRequest) (models.BatchResponseArray, error) {
	processed := s.convertData(ctx, input)
	saved, err := s.Storage.BatchSave(ctx, processed)
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

// DeleteURLs deletes multiple URLs by their short URLs.
func (s *Service) DeleteURLs(ctx context.Context, input models.DeleteURLs) error {
	err := s.Storage.DeleteURLs(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete URLs: %w", err)
	}

	return nil
}

// GetStats gets statistics of saved urls and users.
func (s *Service) GetStats(ctx context.Context) (models.Stats, error) {
	res, err := s.Storage.ServiceStats(ctx)
	if err != nil {
		return models.Stats{}, fmt.Errorf("failed to get stats: %w", err)
	}

	return res, nil
}

func (s *Service) generateUniqueShortLink(ctx context.Context) string {
	const length = 8
	var uniqString string

	// check if the string is unique
	for {
		uniqStringCandidate := generateRandomString(length)
		_, err := s.Storage.Get(ctx, uniqStringCandidate)
		if err != nil {
			if errors.Is(err, ErrURLNotFound) {
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

// GetUserURLs retrieves all URLs associated with a user.
func (s *Service) GetUserURLs(ctx context.Context) (models.UserURLs, error) {
	data, err := s.Storage.GetByUserID(ctx)
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

// IsSubnetTrusted method checks if the IP allowed.
func (s *Service) IsSubnetTrusted(realIP string) bool {
	ipNet, err := parseCIDR(s.TrustedSubnet)
	if err != nil {
		return false
	}

	if realIP != "" {
		clientIP := net.ParseIP(realIP)
		if clientIP != nil && ipNet.Contains(clientIP) {
			return true
		}
	}

	return false
}

func (s *Service) BuildJWTString() (string, error) {
	const tokenExp = time.Hour * 720
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: uuid.NewString(),
	})
	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to create token string: %w", err)
	}
	return tokenString, nil
}

func (s *Service) GetUserID(tokenString, secretKey string, log *logger.Log) string {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Err("failed parse with claims tokenString: ", err)
		return ""
	}
	if !token.Valid {
		log.Err("Token is not valid: ", token)
		return ""
	}

	return claims.UserID
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

func parseCIDR(cidr string) (*net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	return ipNet, nil
}

// ErrURLNotFound error indicates item was not found.
var ErrURLNotFound = errors.New("url not found")

package storage

import (
	"fmt"
	"shortener/internal/config"
)

type Storage struct {
	URLs map[string]string
}

func (s *Storage) Get(shortLink string) (string, bool) {
	longLink, ok := s.URLs[shortLink]
	return longLink, ok
}

func (s *Storage) Save(shortLink, longLink string) {
	s.URLs[shortLink] = longLink
}

func LoadStorage(cfg *config.Config) (*Storage, error) {
	URLs, err := ReadFileStorage(cfg.Service.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file storage %w", err)
	}
	return &Storage{URLs: URLs}, nil
}

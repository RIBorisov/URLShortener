package storage

import (
	"fmt"
	"sync"

	"shortener/internal/config"
	"shortener/internal/logger"
)

type InMemory struct {
	URLs    map[string]string
	mux     *sync.RWMutex
	counter uint64
}

type InFile struct {
	InMemory
	FilePath string
}

type URLStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

func (s *InMemory) Get(shortLink string) (string, bool) {
	s.mux.RLock()
	longLink, ok := s.URLs[shortLink]
	s.mux.RUnlock()
	return longLink, ok
}

func (s *InMemory) Save(shortLink, longLink string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.URLs[shortLink] = longLink
	s.counter++
}

func (s *InFile) Save(shortLink, longLink string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.URLs[shortLink] = longLink
	err := AppendToFile(s.FilePath, shortLink, longLink, s.counter)
	if err != nil {
		logger.Err("failed append to file", err)
	}
	s.counter++
}

func (s *InFile) Restore() error {
	if s.FilePath != "" {
		mapping, err := ReadFileStorage(s.FilePath)
		if err != nil {
			return fmt.Errorf("failed to restore from file %w", err)
		}
		s.mux.Lock()
		s.URLs = mapping
		s.counter = uint64(len(mapping))
		s.mux.Unlock()
	}
	return nil
}

func NewStorage(cfg *config.Config) (URLStorage, error) {
	if cfg.Service.FileStoragePath == "" {
		return &InMemory{
			URLs: make(map[string]string),
			mux:  &sync.RWMutex{},
		}, nil
	}
	storage := &InFile{
		InMemory: InMemory{
			URLs: make(map[string]string),
			mux:  &sync.RWMutex{},
		},
		FilePath: cfg.Service.FileStoragePath,
	}
	err := storage.Restore()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage: %w", err)
	}
	return storage, nil
}

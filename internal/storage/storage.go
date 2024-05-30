package storage

import (
	"fmt"
	"sync"

	"shortener/internal/config"
	"shortener/internal/logger"
)

type inMemory struct {
	urls    map[string]string
	mux     *sync.RWMutex
	counter uint64
}

type InFile struct {
	inMemory
	filePath string
}

type URLStorage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

func (s *inMemory) Get(shortLink string) (string, bool) {
	s.mux.RLock()
	longLink, ok := s.urls[shortLink]
	s.mux.RUnlock()
	return longLink, ok
}

func (s *inMemory) Save(shortLink, longLink string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.urls[shortLink] = longLink
	s.counter++
}

func (s *InFile) Save(shortLink, longLink string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.urls[shortLink] = longLink
	err := AppendToFile(s.filePath, shortLink, longLink, s.counter)
	if err != nil {
		logger.Err("failed append to file", err)
	}
	s.counter++
}

func (s *InFile) restore() error {
	if s.filePath != "" {
		mapping, err := ReadFileStorage(s.filePath)
		if err != nil {
			return fmt.Errorf("failed to restore from file %w", err)
		}
		s.mux.Lock()
		s.urls = mapping
		s.counter = uint64(len(mapping))
		s.mux.Unlock()
	}
	return nil
}

func NewStorage(cfg *config.Config) (URLStorage, error) {
	if cfg.Service.FileStoragePath == "" {
		return &inMemory{
			urls: make(map[string]string),
			mux:  &sync.RWMutex{},
		}, nil
	}
	storage := &InFile{
		inMemory: inMemory{
			urls: make(map[string]string),
			mux:  &sync.RWMutex{},
		},
		filePath: cfg.Service.FileStoragePath,
	}
	err := storage.restore()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage: %w", err)
	}
	return storage, nil
}

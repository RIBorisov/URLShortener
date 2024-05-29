package storage

import (
	"fmt"
	"shortener/internal/logger"
	"sync"

	"shortener/internal/config"
)

type Storage struct {
	URLs     map[string]string
	FilePath string
	mux      *sync.RWMutex
}

func (s *Storage) Get(shortLink string) (string, bool) {
	s.mux.RLock()
	longLink, ok := s.URLs[shortLink]
	s.mux.RUnlock()
	return longLink, ok
}

func (s *Storage) Save(shortLink, longLink string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.URLs[shortLink] = longLink
	if s.FilePath != "" {
		err := AppendToFile(s.FilePath, shortLink, longLink)
		if err != nil {
			logger.Err("failed append to file", err)
		}

	}
}

func (s *Storage) Restore() error {
	if s.FilePath != "" {
		mapping, err := ReadFileStorage(s.FilePath)
		if err != nil {
			return fmt.Errorf("failed to restore from file %w", err)
		}
		s.mux.Lock()
		s.URLs = mapping
		s.mux.Unlock()
	}
	return nil
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	storage := &Storage{
		URLs:     make(map[string]string),
		FilePath: cfg.Service.FileStoragePath,
		mux:      &sync.RWMutex{},
	}
	err := storage.Restore()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage: %w", err)
	}
	return storage, nil

}

/*
Из конструктора хранилища возвращался объект общего интерфейса
	- Существовало две имплементации харинилища - in-memory и файловая
	- Файловое хранилище было оберткой надо in-memory: то есть встаривало бы in-memory при помощи композиции

	В итоге сервису должно быть безразлично, сохраняем ли мы данные в файл или нет.
*/

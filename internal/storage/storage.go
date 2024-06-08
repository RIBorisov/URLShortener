package storage

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/storage/postgres"
)

type inMemory struct {
	urls    map[string]string
	mux     *sync.RWMutex
	counter uint64
}

type inFile struct {
	inMemory
	filePath string
}

type inDatabase struct {
	DB *postgres.DB
}

type URLStorage interface {
	Get(ctx context.Context, shortLink string) (string, bool)
	Save(ctx context.Context, shortLink, longLink string)
}

type URLRow struct {
	Short string `json:"short"`
	Long  string `json:"long"`
	ID    int    `json:"id"`
}

func (i *inDatabase) Get(ctx context.Context, shortLink string) (string, bool) {
	const stmt = `SELECT * FROM urls WHERE short = $1`

	var row URLRow
	err := i.DB.Pool.QueryRow(ctx, stmt, shortLink).Scan(&row)
	if err != nil {
		return "", false // TODO: переписать интерфейс и методы на возвращение error
	}
	return row.Long, true
}

func (i *inDatabase) Save(ctx context.Context, shortLink, longLink string) {
	const stmt = `INSERT INTO urls (short, long) VALUES ($1, $2)`
	res, err := i.DB.Pool.Exec(ctx, stmt, shortLink, longLink)
	if err != nil {
		logger.Err("failed to insert data", err)
	}
	fmt.Println(res.String())
}

func (s *inMemory) Get(_ context.Context, shortLink string) (string, bool) {
	s.mux.RLock()
	longLink, ok := s.urls[shortLink]
	s.mux.RUnlock()
	return longLink, ok
}

func (s *inMemory) Save(_ context.Context, shortLink, longLink string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.urls[shortLink] = longLink
	s.counter++
}

func (s *inFile) Save(_ context.Context, shortLink, longLink string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.urls[shortLink] = longLink
	err := AppendToFile(s.filePath, shortLink, longLink, s.counter)
	if err != nil {
		logger.Err("failed append to file", err)
	}
	s.counter++
}

func (s *inFile) restore() error {
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

func LoadStorage(ctx context.Context, cfg *config.Config) (URLStorage, error) {
	if cfg.Service.DatabaseDSN != "" {
		db, err := postgres.New(ctx, cfg.Service.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to create database storage: %w", err)
		}
		slog.Info("using database storage")
		return &inDatabase{
			DB: db,
		}, nil
	}

	if cfg.Service.FileStoragePath == "" {
		return &inMemory{
			urls: make(map[string]string),
			mux:  &sync.RWMutex{},
		}, nil
	}
	storage := &inFile{
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

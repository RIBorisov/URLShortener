package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/url"
	"sync"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/storage/postgres"
)

type inMemory struct {
	mux     *sync.Mutex
	cfg     *config.Config
	urls    map[string]string
	counter uint64
}

type inFile struct {
	inMemory
	filePath string
}

type inDatabase struct {
	Pool *postgres.DBPool
	cfg  *config.Config
}
type URLStorage interface {
	Get(ctx context.Context, shortLink string) (string, bool)
	Save(ctx context.Context, shortLink, longLink string) error
	BatchSave(ctx context.Context, input models.BatchIn) (models.BatchOut, error)
	Close() error
	Ping(ctx context.Context) error
}

type URLRow struct {
	Short string `json:"short"`
	Long  string `json:"long"`
	ID    int    `json:"id"`
}

func (d *inDatabase) Get(ctx context.Context, shortLink string) (string, bool) {
	const stmt = `SELECT * FROM urls WHERE short = $1`

	var row URLRow
	err := d.Pool.QueryRow(ctx, stmt, shortLink).Scan(&row.ID, &row.Short, &row.Long)
	if err != nil {
		return "", false // TODO: переписать интерфейс и методы на возвращение error
	}
	return row.Long, true
}

func (d *inDatabase) Save(ctx context.Context, shortLink, longLink string) error {
	// потратил добрых часов 6 на правильную реализацию, но не взлетело.
	// Выглядит костыльно, но работает, прошу совета как это исправить/улучшить.
	const insertStmt = `INSERT INTO urls (short, long) VALUES ($1, $2)`
	const selectStmt = `SELECT short FROM urls WHERE long = $1`
	var existingShortLink string
	_, err := d.Pool.Exec(ctx, insertStmt, shortLink, longLink)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			selectErr := d.Pool.Pool.QueryRow(ctx, selectStmt, longLink).Scan(&existingShortLink)
			if selectErr != nil {
				return fmt.Errorf("failed to select row: %w", selectErr)
			}
			return &DuplicateRecordError{Message: existingShortLink, Err: err}
		}
		return fmt.Errorf("failed to execute row: %w", err)
	}
	return nil
}

func (d *inDatabase) BatchSave(ctx context.Context, input models.BatchIn) (models.BatchOut, error) {
	const stmt = `INSERT INTO urls (short, long) VALUES (@correlationID, @originalURL)`

	// получаем connection через pool
	conn, err := d.Pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connectio: %w", err)
	}

	// стартуем tx через connection
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{IsoLevel: "read committed"})
	if err != nil {
		return nil, fmt.Errorf("failed to begion conn tx: %w", err)
	}

	// создаем и наполняем batch
	batch := pgx.Batch{}
	for _, in := range input {
		args := pgx.NamedArgs{
			"correlationID": in.CorrelationID,
			"originalURL":   in.OriginalURL,
		}
		batch.Queue(stmt, args)
	}

	// отдаем в транзакцию и исполняем батчевый запрос
	// batchResults нельзя закрывать в defer т.к. он должен закрыться до(!) закрытия connection и tx.Commit
	batchResults := tx.SendBatch(ctx, &batch)
	_, batchErr := batchResults.Exec()

	if batchErr != nil {
		if err = tx.Rollback(ctx); err != nil {
			logger.Err("failed to rollback transaction", err)
		}
		return nil, fmt.Errorf("batch res exec: %w", err)
	}

	// закрываем тут т.к. нужно дальше коммитить транзакцию
	if err = batchResults.Close(); err != nil {
		logger.Err("failed to close connection results", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	var resp models.BatchOut
	for _, in := range input {
		shortURL, err := url.JoinPath(d.cfg.Service.BaseURL, "/", in.CorrelationID)
		if err != nil {
			return nil, fmt.Errorf("failed to join url: %w", err)
		}
		resp = append(resp, models.BatchResponse{
			CorrelationID: in.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	return resp, nil
}

func (m *inMemory) Get(_ context.Context, shortLink string) (string, bool) {
	m.mux.Lock()
	defer m.mux.Unlock()
	longLink, ok := m.urls[shortLink]
	return longLink, ok
}

func (m *inMemory) Save(_ context.Context, shortLink, longLink string) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.urls[shortLink] = longLink
	m.counter++
	return nil
}

func (m *inMemory) BatchSave(_ context.Context, input models.BatchIn) (models.BatchOut, error) {
	var result models.BatchOut

	for _, item := range input {
		m.mux.Lock()
		m.urls[item.CorrelationID] = item.OriginalURL
		m.mux.Unlock()
		m.counter++
		result = append(result, models.BatchResponse{
			CorrelationID: item.CorrelationID,
			ShortURL:      item.CorrelationID,
		})
	}
	return result, nil
}

func (f *inFile) Save(_ context.Context, shortLink, longLink string) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	f.urls[shortLink] = longLink
	err := AppendToFile(f.filePath, shortLink, longLink, f.counter)
	if err != nil {
		logger.Err("failed append to file", err)
	}
	f.counter++
	return nil
}

func (f *inFile) BatchSave(_ context.Context, input models.BatchIn) (models.BatchOut, error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	saved, err := BatchAppend(f.filePath, f.cfg.Service.BaseURL, input, f.counter)
	if err != nil {
		return nil, fmt.Errorf("failed append rows to file: %w", err)
	}
	f.counter += uint64(len(saved))
	return saved, nil
}

func (f *inFile) restore() error {
	if f.filePath != "" {
		mapping, err := ReadFileStorage(f.filePath)
		if err != nil {
			return fmt.Errorf("failed to restore from file %w", err)
		}
		f.mux.Lock()
		defer f.mux.Unlock()
		f.urls = mapping
		f.counter = uint64(len(mapping))
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
			Pool: db,
			cfg:  cfg,
		}, nil
	}

	if cfg.Service.FileStoragePath == "" {
		return &inMemory{
			urls: make(map[string]string),
			mux:  &sync.Mutex{},
			cfg:  cfg,
		}, nil
	}
	storage := &inFile{
		inMemory: inMemory{
			urls: make(map[string]string),
			mux:  &sync.Mutex{},
			cfg:  cfg,
		},
		filePath: cfg.Service.FileStoragePath,
	}
	err := storage.restore()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage: %w", err)
	}
	return storage, nil
}

func (d *inDatabase) Close() error {
	d.Pool.Close()
	return nil
}

func (m *inMemory) Close() error {
	return nil
}

func (f *inFile) Close() error {
	return nil
}

type DuplicateRecordError struct {
	Err     error
	Message string
}

func (e *DuplicateRecordError) Error() string {
	return e.Message
}
func (e *DuplicateRecordError) Unwrap() error {
	return e.Err
}

func (d *inDatabase) Ping(ctx context.Context) error {
	return d.Pool.Pool.Ping(ctx)
}
func (f *inFile) Ping(_ context.Context) error {
	return nil
}
func (m *inMemory) Ping(_ context.Context) error {
	return nil
}

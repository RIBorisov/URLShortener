// Package storage contains logic for managing URL records.
package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"sync"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/models"
	"shortener/internal/service"
)

// inMemory represents an in-memory URL storage.
type inMemory struct {
	*logger.Log
	mux     *sync.Mutex
	cfg     *config.Config
	urls    map[string]URLRecord
	counter uint64
}

// inFile represents a file-based URL storage.
type inFile struct {
	inMemory
	filePath string
}

// inDatabase represents a database-based URL storage.
type inDatabase struct {
	*DBStore
	cfg *config.Config
	log *logger.Log
}

// Cleanup removes deleted URLs from the database.
func (d *inDatabase) Cleanup(ctx context.Context) ([]string, error) {
	const stmt = `DELETE FROM urls WHERE is_deleted = TRUE RETURNING id`
	result := make([]string, 0)
	rows, err := d.pool.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed query db: %w", err)
	}
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed scan id from row: %w", err)
		}
		result = append(result, id)
	}

	return result, nil
}

// DeleteURLs marks URLs as deleted in the database.
func (d *inDatabase) DeleteURLs(ctx context.Context, input models.DeleteURLs) error {
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return errGetUserFromContext
	}
	const stmt = `UPDATE urls SET is_deleted = TRUE WHERE short = @short AND user_id = @user_id AND is_deleted = FALSE`

	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: "read committed"})
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			d.log.Err("failed to rollback transaction: ", err)
		}
	}()

	batch := pgx.Batch{}
	for _, in := range input {
		args := pgx.NamedArgs{"short": in, "user_id": userID}
		batch.Queue(stmt, args)
	}
	batchResults := tx.SendBatch(ctx, &batch)
	_, batchErr := batchResults.Exec()
	if batchErr != nil {
		return fmt.Errorf("failed execute batch request: %w", batchErr)
	}

	if err = batchResults.Close(); err != nil {
		return fmt.Errorf("failed to close connection results: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByUserID retrieves URLs for a given user ID from the database.
func (d *inDatabase) GetByUserID(ctx context.Context) ([]models.BaseRow, error) {
	const stmt = `SELECT short, long FROM urls WHERE user_id = $1 AND is_deleted = FALSE`
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return nil, errGetUserFromContext
	}
	rows, err := d.pool.Query(ctx, stmt, userID)
	if err != nil {
		return nil, fmt.Errorf("failed get urls for user_id = %s: %w", userID, err)
	}
	var data []models.BaseRow
	for rows.Next() {
		var short, long string
		if err = rows.Scan(&short, &long); err != nil {
			return nil, fmt.Errorf("failed scan rows into BaseRow: %w", err)
		}
		data = append(data, models.BaseRow{
			Short: short,
			Long:  long,
		})
	}
	return data, nil
}

// Get retrieves a URL by its short link from the database.
func (d *inDatabase) Get(ctx context.Context, shortLink string) (string, error) {
	const stmt = `SELECT long, is_deleted FROM urls WHERE short = $1`

	var (
		long      string
		isDeleted bool
	)
	err := d.pool.QueryRow(ctx, stmt, shortLink).Scan(&long, &isDeleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", service.ErrURLNotFound
		}
		return "", fmt.Errorf("failed get row: %w", err)
	}
	if isDeleted {
		return "", ErrURLDeleted
	}
	return long, nil
}

// Save saves a new URL record to the database.
func (d *inDatabase) Save(ctx context.Context, shortLink, longLink string) error {
	const (
		longConstraint = "idx_long_is_not_deleted"
		selectStmt     = `SELECT short FROM urls WHERE long = $1`
		insertStmt     = `INSERT INTO urls (short, long, user_id) VALUES ($1, $2, $3)`
	)
	var existingShortLink string
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return errGetUserFromContext
	}
	// через транзакцию в этом случае нельзя, т.к. если будет получена ошибка, то
	// все последующие команды не будут до роллбэк/коммита выполняться. Savepoints использовать - тут оверхед
	_, err := d.pool.Exec(ctx, insertStmt, shortLink, longLink, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			if pgErr.ConstraintName == longConstraint {
				selectErr := d.pool.QueryRow(ctx, selectStmt, longLink).Scan(&existingShortLink)
				if selectErr != nil {
					return fmt.Errorf("failed to select row: %w", selectErr)
				}
				return &DuplicateRecordError{Message: existingShortLink, Err: err}
			}
		}
		return fmt.Errorf("failed to execute row: %w", err)
	}

	return nil
}

// BatchSave saves multiple URL records to the database.
func (d *inDatabase) BatchSave(ctx context.Context, input models.BatchArray) (models.BatchArray, error) {
	const stmt = `INSERT INTO urls (short, long, user_id) VALUES (@short, @long, @user_id)`

	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: "read committed"})
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			d.log.Err("failed to rollback transaction: ", err)
		}
	}()
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return nil, errGetUserFromContext
	}
	batch := pgx.Batch{}
	for _, in := range input {
		args := pgx.NamedArgs{
			"short":   in.ShortURL,
			"long":    in.OriginalURL,
			"user_id": userID,
		}
		batch.Queue(stmt, args)
	}

	// отдаем в транзакцию и исполняем батчевый запрос
	// batchResults нельзя закрывать в defer т.к. он должен закрыться до(!) закрытия connection и tx.Commit
	batchResults := tx.SendBatch(ctx, &batch)
	_, batchErr := batchResults.Exec()

	if batchErr != nil {
		return nil, fmt.Errorf("failed execute batch request: %w", batchErr)
	}

	// закрываем тут т.к. нужно дальше коммитить транзакцию
	if err = batchResults.Close(); err != nil {
		return nil, fmt.Errorf("failed to close connection results: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	var resp models.BatchArray
	for _, in := range input {
		shortURL, err := url.JoinPath(d.cfg.App.BaseURL, "/", in.ShortURL)
		if err != nil {
			return nil, fmt.Errorf("failed to join url: %w", err)
		}
		resp = append(resp, models.Batch{
			CorrelationID: in.CorrelationID,
			ShortURL:      shortURL,
			OriginalURL:   in.OriginalURL,
		})
	}

	return resp, nil
}

// ServiceStats returns a counter of saved urls and users.
func (d *inDatabase) ServiceStats(ctx context.Context) (models.Stats, error) {
	const selectStmt = `SELECT
    (SELECT COUNT(DISTINCT user_id) FROM urls) AS users_cnt,
    (SELECT COUNT(DISTINCT long) FROM urls) AS long_cnt`

	resp := models.Stats{}

	err := d.pool.QueryRow(ctx, selectStmt).Scan(&resp.Users, &resp.URLs)
	if err != nil {
		return models.Stats{}, fmt.Errorf("failed to get rows from table: %w", err)
	}

	return resp, nil
}

// Cleanup removes deleted URLs from the in-memory storage.
func (m *inMemory) Cleanup(_ context.Context) ([]string, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	cleaned := make([]string, 0)
	for _, u := range m.urls {
		if u.Deleted {
			cleaned = append(cleaned, u.ShortURL)
			delete(m.urls, u.ShortURL)
		}
	}
	return cleaned, nil
}

// GetByUserID retrieves URLs for a given user ID from the in-memory storage.
func (m *inMemory) GetByUserID(ctx context.Context) ([]models.BaseRow, error) {
	var data []models.BaseRow
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return nil, errGetUserFromContext
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, u := range m.urls {
		if userID == u.UserID && !u.Deleted {
			data = append(data, models.BaseRow{
				Long:  u.OriginalURL,
				Short: u.ShortURL,
			})
		}
	}
	return data, nil
}

// DeleteURLs marks URLs as deleted in the in-memory storage.
func (m *inMemory) DeleteURLs(ctx context.Context, input models.DeleteURLs) error {
	var wg sync.WaitGroup
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return errGetUserFromContext
	}
	for _, short := range input {
		short := short
		wg.Add(1)
		go func() {
			defer wg.Done()
			u, ok := m.urls[short]
			if !ok {
				m.Log.Err("url not found", short)
				return
			}

			if u.UserID == userID && !u.Deleted {
				m.urls[short] = URLRecord{
					OriginalURL: u.OriginalURL,
					ShortURL:    u.ShortURL,
					UUID:        u.UUID,
					UserID:      u.UserID,
					Deleted:     true,
				}
				m.Log.Debug("deleted url", "short", u.ShortURL)
			}
		}()
	}
	wg.Wait()

	return nil
}

// Get retrieves a URL by its short link from the in-memory storage.
func (m *inMemory) Get(_ context.Context, shortLink string) (string, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	longLink, ok := m.urls[shortLink]
	if ok {
		return longLink.OriginalURL, nil
	}
	return "", service.ErrURLNotFound
}

// Save saves a new URL record to the in-memory storage.
func (m *inMemory) Save(ctx context.Context, shortLink, longLink string) error {
	m.mux.Lock()
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return errGetUserFromContext
	}
	defer m.mux.Unlock()
	m.urls[shortLink] = URLRecord{
		UUID:        strconv.FormatUint(m.counter, 10),
		OriginalURL: longLink,
		ShortURL:    shortLink,
		UserID:      userID,
		Deleted:     false,
	}
	m.counter++
	return nil
}

// BatchSave saves multiple URL records to the in-memory storage.
func (m *inMemory) BatchSave(ctx context.Context, input models.BatchArray) (models.BatchArray, error) {
	var result models.BatchArray
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return nil, errGetUserFromContext
	}
	for _, item := range input {
		m.mux.Lock()
		m.urls[item.ShortURL] = URLRecord{OriginalURL: item.OriginalURL, UUID: item.CorrelationID, UserID: userID}
		m.mux.Unlock()
		m.counter++
		result = append(result, models.Batch{
			CorrelationID: item.CorrelationID,
			ShortURL:      item.CorrelationID,
			OriginalURL:   item.OriginalURL,
		})
	}
	return result, nil
}

// ServiceStats returns a counter of saved urls and users.
func (m *inMemory) ServiceStats(_ context.Context) (models.Stats, error) {
	result := models.Stats{
		URLs:  int(m.counter),
		Users: getUniqUsers(m.urls),
	}
	return result, nil
}

func getUniqUsers(m map[string]URLRecord) int {
	uniqUsers := make(map[string]struct{})
	for _, urlRecord := range m {
		uniqUsers[urlRecord.UserID] = struct{}{}
	}
	return len(uniqUsers)
}

// Cleanup removes deleted URLs from the file-based storage.
func (f *inFile) Cleanup(_ context.Context) ([]string, error) {
	urls := make([]URLRecord, 0)
	f.mux.Lock()
	for _, u := range f.inMemory.urls {
		if !u.Deleted {
			urls = append(urls, URLRecord{
				UUID:        u.UUID,
				OriginalURL: u.OriginalURL,
				ShortURL:    u.ShortURL,
				UserID:      u.UserID,
				Deleted:     u.Deleted,
			})
		}
	}
	f.mux.Unlock()

	if err := BatchUpdate(f.filePath, urls); err != nil {
		return nil, fmt.Errorf("failed batch update: %w", err)
	}
	result := make([]string, 0, len(urls))
	for _, u := range urls {
		result = append(result, u.ShortURL)
	}

	return result, nil
}

// Save saves a new URL record to the file-based storage.
func (f *inFile) Save(ctx context.Context, shortLink, longLink string) error {
	f.mux.Lock()
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return errGetUserFromContext
	}
	defer f.mux.Unlock()
	urlRecord := URLRecord{
		UUID:        strconv.FormatUint(f.counter+1, 10),
		OriginalURL: longLink,
		UserID:      userID,
		ShortURL:    shortLink,
	}
	f.urls[shortLink] = urlRecord

	err := AppendToFile(f.Log, f.filePath, urlRecord)
	if err != nil {
		return fmt.Errorf("failed append to file: %w", err)
	}
	f.counter++
	return nil
}

// BatchSave saves multiple URL records to the file-based storage.
func (f *inFile) BatchSave(ctx context.Context, input models.BatchArray) (models.BatchArray, error) {
	f.mux.Lock()
	userID, ok := ctx.Value(models.CtxUserIDKey).(string)
	if !ok {
		return nil, errGetUserFromContext
	}
	defer f.mux.Unlock()
	saved, err := BatchAppend(f.Log, f.filePath, f.cfg.App.BaseURL, userID, input, f.counter)
	if err != nil {
		return nil, fmt.Errorf("failed append rows to file: %w", err)
	}
	f.counter += uint64(len(saved))
	return saved, nil
}

// DeleteURLs marks URLs as deleted in the file-based storage.
func (f *inFile) DeleteURLs(ctx context.Context, input models.DeleteURLs) error {
	err := f.inMemory.DeleteURLs(ctx, input)
	if err != nil {
		return fmt.Errorf("failed delete user urls: %w", err)
	}
	urls := make([]URLRecord, 0, len(input))
	f.mux.Lock()
	defer f.mux.Unlock()

	for _, u := range f.inMemory.urls {
		urls = append(urls, URLRecord{
			UUID:        u.UUID,
			OriginalURL: u.OriginalURL,
			ShortURL:    u.ShortURL,
			UserID:      u.UserID,
			Deleted:     u.Deleted,
		})
	}

	if err = BatchUpdate(f.filePath, urls); err != nil {
		return fmt.Errorf("failed batch update: %w", err)
	}

	return nil
}

// restore restores the file-based storage from a file.
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

// ServiceStats returns a counter of saved urls and users.
func (f *inFile) ServiceStats(ctx context.Context) (models.Stats, error) {
	res, err := f.inMemory.ServiceStats(ctx)
	if err != nil {
		return models.Stats{}, fmt.Errorf("failed to get stats: %w", err)
	}
	return res, nil
}

// LoadStorage loads the appropriate URL storage based on the configuration.
func LoadStorage(ctx context.Context, cfg *config.Config, log *logger.Log) (service.URLStorage, error) {
	if cfg.App.DatabaseDSN != "" {
		db, err := New(ctx, cfg.App.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to create database storage: %w", err)
		}
		log.Info("using database storage..")
		return &inDatabase{db, cfg, log}, nil
	}

	if cfg.App.FileStoragePath == "" {
		log.Info("using memory storage..")
		return &inMemory{
			urls: make(map[string]URLRecord),
			mux:  &sync.Mutex{},
			cfg:  cfg,
			Log:  log,
		}, nil
	}
	storage := &inFile{
		inMemory: inMemory{
			urls: make(map[string]URLRecord),
			mux:  &sync.Mutex{},
			cfg:  cfg,
			Log:  log,
		},
		filePath: cfg.App.FileStoragePath,
	}
	err := storage.restore()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage: %w", err)
	}
	log.Debug("using file storage..")

	return storage, nil
}

// Close closes the database connection.
func (d *inDatabase) Close() error {
	d.pool.Close()
	return nil
}

// Close does nothing for in-memory storage.
func (m *inMemory) Close() error {
	return nil
}

// Close does nothing for in-file storage.
func (f *inFile) Close() error {
	return nil
}

// DuplicateRecordError ...
type DuplicateRecordError struct {
	Err     error
	Message string
}

// Error ...
func (e *DuplicateRecordError) Error() string {
	return e.Message
}

// Unwrap ...
func (e *DuplicateRecordError) Unwrap() error {
	return e.Err
}

// Ping ...
func (d *inDatabase) Ping(ctx context.Context) error {
	return d.pool.Ping(ctx)
}

// Ping ...
func (f *inFile) Ping(_ context.Context) error {
	return nil
}

// Ping ...
func (m *inMemory) Ping(_ context.Context) error {
	return nil
}

// ErrURLDeleted ...
var (
	ErrURLDeleted         = errors.New("url has been deleted")
	errGetUserFromContext = errors.New("failed get user from context")
)

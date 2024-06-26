package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/models"
)

type inMemory struct {
	*logger.Log
	mux     *sync.Mutex
	cfg     *config.Config
	urls    map[string]URLRecord
	counter uint64
}

type inFile struct {
	inMemory
	filePath string
}

type inDatabase struct {
	*DBStore
	cfg *config.Config
	log *logger.Log
}

type URLStorage interface {
	Close() error
	Ping(ctx context.Context) error
	Get(ctx context.Context, shortLink string) (string, error)
	Save(ctx context.Context, shortLink, longLink string, user *models.User) error
	BatchSave(ctx context.Context, input models.BatchArray, user *models.User) (models.BatchArray, error)
	GetByUserID(ctx context.Context, user *models.User) ([]models.BaseRow, error)
	DeleteURLs(ctx context.Context, input models.DeleteURLs, user *models.User) error
}

func (d *inDatabase) DeleteURLs(ctx context.Context, input models.DeleteURLs, user *models.User) error {
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
		args := pgx.NamedArgs{"short": in, "user_id": user.ID}
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

func (d *inDatabase) GetByUserID(ctx context.Context, user *models.User) ([]models.BaseRow, error) {
	const stmt = `SELECT short, long FROM urls WHERE user_id = $1 AND is_deleted = FALSE`
	rows, err := d.pool.Query(ctx, stmt, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed get urls for user_id = %s: %w", user.ID, err)
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

func (d *inDatabase) Get(ctx context.Context, shortLink string) (string, error) {
	const stmt = `SELECT long, is_deleted FROM urls WHERE short = $1`

	var (
		long      string
		isDeleted bool
	)
	err := d.pool.QueryRow(ctx, stmt, shortLink).Scan(&long, &isDeleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrURLNotFound
		}
		return "", fmt.Errorf("failed get row: %w", err)
	}
	if isDeleted {
		return "", ErrURLDeleted
	}
	return long, nil
}

func (d *inDatabase) Save(ctx context.Context, shortLink, longLink string, user *models.User) error {
	const (
		longConstraint = "idx_long_url"
		selectStmt     = `SELECT short FROM urls WHERE long = $1`
		insertStmt     = `INSERT INTO urls (short, long, user_id) VALUES ($1, $2, $3)`
	)
	var existingShortLink string
	// через транзакцию в этом случае нельзя, т.к. если будет получена ошибка, то
	// все последующие команды не будут до роллбэк/коммита выполняться. Savepoints использовать - тут оверхед
	_, err := d.pool.Exec(ctx, insertStmt, shortLink, longLink, user.ID)
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

func (d *inDatabase) BatchSave(
	ctx context.Context,
	input models.BatchArray,
	user *models.User,
) (models.BatchArray, error) {
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

	batch := pgx.Batch{}
	for _, in := range input {
		args := pgx.NamedArgs{
			"short":   in.ShortURL,
			"long":    in.OriginalURL,
			"user_id": user.ID,
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
		shortURL, err := url.JoinPath(d.cfg.Service.BaseURL, "/", in.ShortURL)
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

func (m *inMemory) GetByUserID(_ context.Context, user *models.User) ([]models.BaseRow, error) {
	var data []models.BaseRow
	for _, u := range m.urls {
		m.mux.Lock()
		if user.ID == u.UserID && u.Deleted == false {
			data = append(data, models.BaseRow{
				Long:  u.OriginalURL,
				Short: u.ShortURL,
			})
		}
		m.mux.Unlock()
	}
	return data, nil
}

func (m *inMemory) DeleteURLs(_ context.Context, input models.DeleteURLs, user *models.User) error {
	var wg sync.WaitGroup
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

			if u.UserID == user.ID && u.Deleted == false {
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

func (m *inMemory) Get(_ context.Context, shortLink string) (string, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	longLink, ok := m.urls[shortLink]
	if ok {
		return longLink.OriginalURL, nil
	}
	return "", ErrURLNotFound
}

func (m *inMemory) Save(_ context.Context, shortLink, longLink string, user *models.User) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.urls[shortLink] = URLRecord{
		OriginalURL: longLink,
		ShortURL:    shortLink,
		UserID:      user.ID,
		Deleted:     false,
	}
	m.counter++
	return nil
}

func (m *inMemory) BatchSave(_ context.Context, input models.BatchArray, user *models.User) (models.BatchArray, error) {
	var result models.BatchArray

	for _, item := range input {
		m.mux.Lock()
		m.urls[item.ShortURL] = URLRecord{OriginalURL: item.OriginalURL, UUID: item.CorrelationID, UserID: user.ID}
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

func (f *inFile) Save(_ context.Context, shortLink, longLink string, user *models.User) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	f.urls[shortLink] = URLRecord{OriginalURL: longLink, UserID: user.ID, ShortURL: shortLink}
	err := AppendToFile(f.Log, f.filePath, shortLink, longLink, f.counter, user)
	if err != nil {
		return fmt.Errorf("failed append to file: %w", err)
	}
	f.counter++
	return nil
}

func (f *inFile) BatchSave(_ context.Context, input models.BatchArray, user *models.User) (models.BatchArray, error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	saved, err := BatchAppend(f.Log, f.filePath, f.cfg.Service.BaseURL, input, f.counter, user)
	if err != nil {
		return nil, fmt.Errorf("failed append rows to file: %w", err)
	}
	f.counter += uint64(len(saved))
	return saved, nil
}

func (f *inFile) DeleteURLs(ctx context.Context, input models.DeleteURLs, user *models.User) error {
	err := f.inMemory.DeleteURLs(ctx, input, user)
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

func LoadStorage(ctx context.Context, cfg *config.Config, log *logger.Log) (URLStorage, error) {
	if cfg.Service.DatabaseDSN != "" {
		db, err := New(ctx, cfg.Service.DatabaseDSN, log)
		if err != nil {
			return nil, fmt.Errorf("failed to create database storage: %w", err)
		}
		log.Info("using database storage..")
		return &inDatabase{db, cfg, log}, nil
	}

	if cfg.Service.FileStoragePath == "" {
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
		filePath: cfg.Service.FileStoragePath,
	}
	err := storage.restore()
	if err != nil {
		return nil, fmt.Errorf("failed to build storage: %w", err)
	}
	log.Info("using file storage..")

	return storage, nil
}

func (d *inDatabase) Close() error {
	d.pool.Close()
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
	return d.pool.Ping(ctx)
}
func (f *inFile) Ping(_ context.Context) error {
	return nil
}
func (m *inMemory) Ping(_ context.Context) error {
	return nil
}

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLDeleted  = errors.New("url has been deleted")
)

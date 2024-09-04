package storage

import (
	"context"
	"errors"
	"os"
	"path"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"shortener/internal/config"
	"shortener/internal/logger"
	"shortener/internal/models"
)

const (
	baseLongURL  = "https://example.com"
	baseShortURL = "short"
)

func TestClose_InMemory(t *testing.T) {
	memStorage := &inMemory{}
	if err := memStorage.Close(); err != nil {
		assert.NoError(t, err)
	}
}

func TestClose_InFile(t *testing.T) {
	fileStorage := &inFile{}
	if err := fileStorage.Close(); err != nil {
		assert.NoError(t, err)
	}
}

func TestDuplicateRecordError_Error(t *testing.T) {
	err := errors.New("duplicate record")
	dupErr := &DuplicateRecordError{Err: err, Message: "duplicate found"}
	if dupErr.Error() != "duplicate found" {
		t.Errorf("expected 'duplicate found', got '%s'", dupErr.Error())
	}
}

func TestDuplicateRecordError_Unwrap(t *testing.T) {
	err := errors.New("duplicate record")
	dupErr := &DuplicateRecordError{Err: err, Message: "duplicate found"}
	if !errors.Is(dupErr.Unwrap(), err) {
		t.Errorf("expected %v, got %v", err, dupErr.Unwrap())
	}
}

func TestPing_InMemory(t *testing.T) {
	memStorage := &inMemory{}
	if err := memStorage.Ping(context.Background()); err != nil {
		assert.NoError(t, err)
	}
}

func TestPing_InFile(t *testing.T) {
	fileStorage := &inFile{}
	if err := fileStorage.Ping(context.Background()); err != nil {
		assert.NoError(t, err)
	}
}

func TestSave_InMemory(t *testing.T) {
	log := &logger.Log{}
	log.Initialize("INFO")
	memStorage := &inMemory{
		mux:     &sync.Mutex{},
		counter: 0,
		Log:     log,
		urls:    make(map[string]URLRecord),
	}
	ctx := context.WithValue(context.Background(), models.CtxUserIDKey, "user_id")
	if err := memStorage.Save(ctx, baseShortURL, baseLongURL); err != nil {
		assert.NoError(t, err)
	}

	longLink, err := memStorage.Get(ctx, baseShortURL)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, baseLongURL, longLink)
}

func TestGet_InMemory(t *testing.T) {
	log := &logger.Log{}
	log.Initialize("INFO")
	memStorage := &inMemory{
		mux:     &sync.Mutex{},
		counter: 0,
		Log:     log,
		urls:    make(map[string]URLRecord),
	}
	if err := memStorage.Save(
		context.WithValue(context.Background(), models.CtxUserIDKey, "user_id"),
		baseShortURL,
		baseLongURL,
	); err != nil {
		assert.NoError(t, err)
	}

	longLink, err := memStorage.Get(context.Background(), baseShortURL)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, baseLongURL, longLink)
}

func TestBatchSave_InMemory(t *testing.T) {
	log := &logger.Log{}
	log.Initialize("INFO")
	memStorage := &inMemory{
		mux:     &sync.Mutex{},
		counter: 0,
		Log:     log,
		urls:    make(map[string]URLRecord),
	}
	ctx := context.WithValue(context.Background(), models.CtxUserIDKey, "user_id")

	batchInput := models.BatchArray{
		{ShortURL: "short1", OriginalURL: "https://example.com/1", CorrelationID: "1"},
		{ShortURL: "short2", OriginalURL: "https://example.com/2", CorrelationID: "2"},
	}

	saved, err := memStorage.BatchSave(ctx, batchInput)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, len(batchInput), len(saved))
}

func TestInMemoryCleanup(t *testing.T) {
	const (
		short1 = "short1"
		short2 = "short2"
		user1  = "user1"
		user2  = "user2"
	)
	log := &logger.Log{}
	log.Initialize("INFO")
	mem := &inMemory{
		Log:  log,
		mux:  &sync.Mutex{},
		cfg:  &config.Config{},
		urls: make(map[string]URLRecord),
	}

	mem.urls["short1"] = URLRecord{
		UUID:        "1",
		OriginalURL: "https://example1.com",
		ShortURL:    short1,
		UserID:      user1,
		Deleted:     true,
	}
	mem.urls["short2"] = URLRecord{
		UUID:        "2",
		OriginalURL: "https://example2.com",
		ShortURL:    short2,
		UserID:      user2,
		Deleted:     false,
	}

	cleaned, err := mem.Cleanup(context.Background())
	assert.NoError(t, err)

	expected := []string{short1}
	assert.Equal(t, len(expected), len(cleaned))

	if _, ok := mem.urls[short1]; ok {
		t.Errorf("expected URL %s to be deleted, but it still exists", short1)
	}
	if _, ok := mem.urls["short2"]; !ok {
		t.Errorf("expected URL %s to still exist, but it was deleted", short2)
	}
}

func TestInFileDeleteURLs(t *testing.T) {
	log := &logger.Log{}
	log.Initialize("INFO")
	tmpDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().Unix(), 10))
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)
	tmpFile, err := os.CreateTemp(tmpDir, filename)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, tmpFile.Close())
	}()
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			assert.NoError(t, err)
		}
	}()
	ctx := context.WithValue(context.Background(), models.CtxUserIDKey, "user1")

	file := &inFile{
		inMemory: inMemory{
			Log:  log,
			mux:  &sync.Mutex{},
			cfg:  &config.Config{},
			urls: make(map[string]URLRecord),
		},
		filePath: tmpFile.Name(),
	}

	file.inMemory.urls["short1"] = URLRecord{
		UUID:        "1",
		OriginalURL: "https://example1.com",
		ShortURL:    "short1",
		UserID:      "user1",
		Deleted:     false,
	}
	file.inMemory.urls["short2"] = URLRecord{
		UUID:        "2",
		OriginalURL: "https://example2.com",
		ShortURL:    "short2",
		UserID:      "user2",
		Deleted:     false,
	}

	input := models.DeleteURLs{"short1"}
	err = file.DeleteURLs(ctx, input)
	assert.NoError(t, err)

	if _, ok := file.inMemory.urls["short1"]; ok && !file.inMemory.urls["short1"].Deleted {
		t.Errorf("expected URL short1 to be deleted, but it still exists")
	}
	if _, ok := file.inMemory.urls["short2"]; !ok || file.inMemory.urls["short2"].Deleted {
		t.Errorf("expected URL short2 to still exist and not be deleted, but it was deleted")
	}

	urls, err := ReadFileStorage(file.filePath)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(urls))
	assert.True(t, urls["short1"].Deleted)
	assert.False(t, urls["short2"].Deleted)
}

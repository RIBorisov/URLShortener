package storage

import (
	"encoding/json"
	"os"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"shortener/internal/logger"
	"shortener/internal/models"
)

const filename = "tmp_file.txt"

func TestNewConsumer(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().Unix(), 10))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	filePath := path.Join(tmpDir, filename)
	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			assert.NoError(t, err)
		}
	}()

	consumer, err := NewConsumer(filePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer func() {
		if err = consumer.file.Close(); err != nil {
			assert.NoError(t, err)
		}
	}()
	assert.NotNil(t, consumer.file)
}

func TestReadFileStorage(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().Unix(), 10))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		assert.NoError(t, err)
	}
	filePath := path.Join(tmpDir, filename)
	defer func() {
		err := os.RemoveAll(tmpDir)
		assert.NoError(t, err)
	}()

	records := []URLRecord{
		{UUID: "1", OriginalURL: "https://example.com/1", ShortURL: "short1", UserID: "user1"},
		{UUID: "2", OriginalURL: "https://example.com/2", ShortURL: "short2", UserID: "user2"},
	}

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		assert.NoError(t, err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			assert.NoError(t, err)
		}
	}()

	for _, record := range records {
		data, err := json.Marshal(record)
		if err != nil {
			assert.NoError(t, err)
		}
		data = append(data, '\n')
		if _, err = file.Write(data); err != nil {
			assert.NoError(t, err)
		}
	}

	urls, err := ReadFileStorage(filePath)
	if err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, len(urls), len(records))
}

func TestAppendToFile(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().Unix(), 10))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		assert.NoError(t, err)
	}
	filePath := path.Join(tmpDir, filename)
	defer func() {
		err := os.RemoveAll(tmpDir)
		assert.NoError(t, err)
	}()

	log := &logger.Log{}
	log.Initialize("INFO")
	record := URLRecord{
		UUID:        "1",
		OriginalURL: "https://example.com/1",
		ShortURL:    "short1",
		UserID:      "user1",
	}

	if err := AppendToFile(log, filePath, record); err != nil {
		assert.NoError(t, err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		assert.NoError(t, err)
	}

	var readRecord URLRecord
	if err := json.Unmarshal(content[:len(content)-1], &readRecord); err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, readRecord, record)
}

func TestBatchAppend(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().Unix(), 10))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		assert.NoError(t, err)
	}
	filePath := path.Join(tmpDir, filename)
	defer func() {
		err := os.RemoveAll(tmpDir)
		assert.NoError(t, err)
	}()

	log := &logger.Log{}
	log.Initialize("INFO")
	baseURL := "http://short.url"

	input := models.BatchArray{
		{OriginalURL: "http://example.com/1", CorrelationID: "short1"},
		{OriginalURL: "http://example.com/2", CorrelationID: "short2"},
	}

	saved, err := BatchAppend(log, filePath, baseURL, "user1", input, 0)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, 2, len(saved))
}

func TestBatchUpdate(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), strconv.FormatInt(time.Now().Unix(), 10))
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		assert.NoError(t, err)
	}

	records := []URLRecord{
		{UUID: "1", OriginalURL: "https://example.com/3", ShortURL: "short3", UserID: "user3"},
		{UUID: "2", OriginalURL: "https://example.com/4", ShortURL: "short4", UserID: "user4"},
	}

	file, err := os.CreateTemp(tmpDir, filename)
	assert.NoError(t, err)
	defer func() {
		err = file.Close()
		if err != nil {
			assert.NoError(t, err)
		}
	}()

	if err = BatchUpdate(file.Name(), records); err != nil {
		assert.NoError(t, err)
	}

	urls, err := ReadFileStorage(file.Name())
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, len(urls), len(records))
}

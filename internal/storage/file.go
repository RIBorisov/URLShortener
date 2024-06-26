package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"shortener/internal/logger"
	"shortener/internal/models"
)

type URLRecord struct {
	UUID        string `json:"uuid"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	UserID      string `json:"user_id"`
	Deleted     bool   `json:"is_deleted"`
}
type Consumer struct {
	file   *os.File
	reader *bufio.Scanner
}

func prepareDir(filename string) error {
	dir := filepath.Dir(filename)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0666)
		if err != nil {
			return fmt.Errorf("failed to create directories: %w", err)
		}
	}
	return nil
}

func NewConsumer(filename string) (*Consumer, error) {
	err := prepareDir(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare direcory: %w", err)
	}
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	return &Consumer{
		file:   file,
		reader: bufio.NewScanner(file),
	}, nil
}

func ReadFileStorage(filename string) (map[string]URLRecord, error) {
	c, err := NewConsumer(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create new consumer: %w", err)
	}
	var urlRecord URLRecord
	var URLs = map[string]URLRecord{}

	for c.reader.Scan() {
		row := c.reader.Text()
		err = json.Unmarshal([]byte(row), &urlRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal row: %w", err)
		}
		URLs[urlRecord.ShortURL] = URLRecord{
			UUID:        urlRecord.UUID,
			OriginalURL: urlRecord.OriginalURL,
			ShortURL:    urlRecord.ShortURL,
			UserID:      urlRecord.UserID,
		}
	}

	return URLs, nil
}

func AppendToFile(log *logger.Log, filename, short, long string, uuid uint64, user *models.User) error {
	urlRecord := URLRecord{
		UUID:        strconv.FormatUint(uuid+1, 10),
		OriginalURL: long,
		ShortURL:    short,
		UserID:      user.ID,
		Deleted:     false,
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Err("failed to close file", err)
		}
	}()
	data, err := json.Marshal(&urlRecord)
	if err != nil {
		return fmt.Errorf("failed to marshal url record %w", err)
	}

	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed write to file %w", err)
	}

	return nil
}

func BatchAppend(
	log *logger.Log,
	filename,
	baseURL string,
	input models.BatchArray,
	counter uint64,
	user *models.User,
) (models.BatchArray, error) {
	var saved models.BatchArray
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Err("failed to close file: ", err)
		}
	}()
	for _, item := range input {
		var row = URLRecord{
			UUID:        strconv.FormatUint(counter+1, 10),
			OriginalURL: item.OriginalURL,
			ShortURL:    item.CorrelationID,
			UserID:      user.ID,
			Deleted:     false,
		}
		data, err := json.Marshal(&row)
		if err != nil {
			return nil, fmt.Errorf("failed marshal row: %w", err)
		}
		data = append(data, '\n')
		_, err = file.Write(data)
		if err != nil {
			return nil, fmt.Errorf("failed write batch into file: %w", err)
		}
		counter++
		shortURL, err := url.JoinPath(baseURL, "/", item.CorrelationID)
		if err != nil {
			return nil, fmt.Errorf("failed to build short url: %w", err)
		}
		saved = append(saved, models.Batch{
			CorrelationID: item.CorrelationID,
			ShortURL:      shortURL,
			OriginalURL:   item.OriginalURL,
		})
	}

	return saved, nil
}
func BatchUpdate(filename string, input []URLRecord) error {
	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file %w", err)
	}
	result := make([]byte, 0)
	for _, in := range input {
		record := URLRecord{
			UUID:        in.UUID,
			OriginalURL: in.OriginalURL,
			ShortURL:    in.ShortURL,
			UserID:      in.UserID,
			Deleted:     in.Deleted,
		}
		data, err := json.Marshal(&record)
		if err != nil {
			return fmt.Errorf("failed marshal data: %w", err)
		}
		data = append(data, '\n')
		result = append(result, data...)
	}
	if _, err = file.Write(result); err != nil {
		return fmt.Errorf("failed write line into file: %w", err)
	}

	return nil
}

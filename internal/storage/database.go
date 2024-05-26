package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"shortener/internal/logger"
)

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func ReadFileStorage(filename string) (map[string]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %w", err)
	}

	var URLs = map[string]string{}
	rows := strings.Split(string(data), "\n")
	for _, row := range rows {
		if row == "" {
			continue
		}
		var urlRecord URLRecord
		err = json.Unmarshal([]byte(row), &urlRecord)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal row %w", err)
		}
		URLs[urlRecord.ShortURL] = urlRecord.OriginalURL
	}

	return URLs, nil
}

func AppendToFile(filename, short, long string) error {
	uuid, err := GenerateNextUUID(filename)
	if err != nil {
		return fmt.Errorf("failed to generate next uuid %w", err)
	}
	urlRecord := URLRecord{
		UUID:        uuid,
		ShortURL:    short,
		OriginalURL: long,
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file %w", err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			logger.Err("failed to close file %w", err)
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

func GenerateNextUUID(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %w", err)
	}
	rows := strings.Split(string(data), "\n")

	return strconv.Itoa(len(rows)), nil
}

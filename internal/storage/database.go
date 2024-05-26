package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"shortener/internal/logger"
)

type URLRecord struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
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

func ReadFileStorage(filename string) (map[string]string, error) {
	consumer, err := NewConsumer(filename)
	if err != nil {
		logger.Err("failed to create new consumer: ", err)
		return nil, fmt.Errorf("failed to create new consumer: %w", err)
	}
	data := consumer.reader.Bytes()
	//data, err := os.ReadFile(filename)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to read file %w", err)
	//}

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

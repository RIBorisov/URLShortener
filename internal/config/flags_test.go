package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f = ServiceConfig{}

	// Задаем аргументы командной строки
	os.Args = []string{"cmd", "-a", "127.0.0.1:9090", "-b", "http://127.0.0.1:9090", "-f", "/path/to/storage", "-d", "user:pass@tcp(localhost:3306)/dbname"}

	parsed := parseFlags()

	assert.Equal(t, "127.0.0.1:9090", parsed.ServerAddress)
	assert.Equal(t, "http://127.0.0.1:9090", parsed.BaseURL)
	assert.Equal(t, "/path/to/storage", parsed.FileStoragePath)
	assert.Equal(t, "user:pass@tcp(localhost:3306)/dbname", parsed.DatabaseDSN)

	newConfig := parseFlags()

	assert.Equal(t, parsed, newConfig) // Проверяем, что результат не изменился в случае повторого парсинга
}

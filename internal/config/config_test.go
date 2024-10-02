package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Positive #1",
			want: &Config{
				App: AppConfig{
					ServerAddress:    ":8080",
					ServerAddressTLS: ":8443",
					BaseURL:          "base-url.com",
					FileStoragePath:  "/tmp/file.json",
					DatabaseDSN:      "postgresql://user:password@db-host:5432/db-name?sslmode=false",
					EnableHTTPS:      false,
				},
				Service: ServiceConfig{
					SecretKey:                 "super",
					BackgroundCleanup:         true,
					BackgroundCleanupInterval: time.Duration(60000000000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("SERVER_ADDRESS", ":8080")
			t.Setenv("SERVER_ADDRESS_TLS", ":8443")
			t.Setenv("BASE_URL", "base-url.com")
			t.Setenv("FILE_STORAGE_PATH", "/tmp/file.json")
			t.Setenv("DATABASE_DSN", "postgresql://user:password@db-host:5432/db-name?sslmode=false")
			t.Setenv("SECRET_KEY", "super")
			t.Setenv("BACKGROUND_CLEANUP", "0")
			assert.Equal(t, tt.want, LoadConfig())
		})
	}
}

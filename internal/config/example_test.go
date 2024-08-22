package config

import (
	"fmt"
	"time"
)

// ExampleConfig demonstrates how to load and use the configuration.
func ExampleConfig() {
	// Load the configuration.
	cfg := Config{
		Service: ServiceConfig{
			ServerAddress:             ":8080",
			BaseURL:                   "http://localhost:8080",
			FileStoragePath:           "/tmp/file-storage-path/file.json",
			DatabaseDSN:               "postgresql://admin:password@localhost:5432/shortener?sslmode=disable",
			SecretKey:                 "super-secret-key",
			BackgroundCleanup:         true,
			BackgroundCleanupInterval: 10 * time.Second,
		},
	}
	// Configure output
	output := fmt.Sprintf(
		`{
	ServerAddress: %q,
	BaseURL: %q,
	FileStoragePath: %q,
	DatabaseDSN: %q,
	SecretKey: %q,
	BackgroundCleanup: %t,
	BackgroundCleanupInterval: %v
}`,
		cfg.Service.ServerAddress,
		cfg.Service.BaseURL,
		cfg.Service.FileStoragePath,
		cfg.Service.DatabaseDSN,
		cfg.Service.SecretKey,
		cfg.Service.BackgroundCleanup,
		cfg.Service.BackgroundCleanupInterval,
	)

	// Print the loaded configuration.
	fmt.Println(output)

	// Output:
	// {
	// 	ServerAddress: ":8080",
	// 	BaseURL: "http://localhost:8080",
	// 	FileStoragePath: "/tmp/file-storage-path/file.json",
	// 	DatabaseDSN: "postgresql://admin:password@localhost:5432/shortener?sslmode=disable",
	// 	SecretKey: "super-secret-key",
	// 	BackgroundCleanup: true,
	// 	BackgroundCleanupInterval: 10s
	// }
}

package logger

import (
	"log/slog"
)

func Initialize() *slog.Logger {
	//return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return slog.Default()
}

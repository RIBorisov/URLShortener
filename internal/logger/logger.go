package logger

import (
	"fmt"
	"log/slog"
)

func Initialize() *slog.Logger {
	// return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	return slog.Default()
}

// Err обертка для логирования ошибки
func Err(message string, value interface{}) {
	slog.Error(message, slog.String("err", fmt.Sprintf("%v", value)))
}

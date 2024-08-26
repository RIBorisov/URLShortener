package logger

import (
	"fmt"
	"log/slog"
	"os"
)

// Log struct.
type Log struct {
	Logger *slog.Logger
}

// Initialize logger.
func (l *Log) Initialize(level string) *slog.Logger {
	const (
		info    = "INFO"
		debug   = "DEBUG"
		warning = "WARNING"
		errorL  = "ERROR"
	)
	var logLevel slog.Level
	switch level {
	case info:
		logLevel = slog.LevelInfo
	case debug:
		logLevel = slog.LevelDebug
	case warning:
		logLevel = slog.LevelWarn
	case errorL:
		logLevel = slog.LevelError
	}
	l.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	return l.Logger
}

// Fatal event.
func (l *Log) Fatal(v ...any) {
	l.Logger.Error(fmt.Sprint(v...))
	os.Exit(1)
}

// Err event.
func (l *Log) Err(message string, value interface{}) {
	slog.Error(message, slog.String("err", fmt.Sprintf("%v", value)))
}

// Info event.
func (l *Log) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

// Debug event.
func (l *Log) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

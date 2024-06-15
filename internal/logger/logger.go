package logger

import (
	"fmt"
	"log/slog"
	"os"
)

type Log struct {
	Logger *slog.Logger
}

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

func (l *Log) Fatal(v ...any) {
	l.Logger.Error(fmt.Sprint(v...))
	os.Exit(1)
}

func (l *Log) Err(message string, value interface{}) {
	slog.Error(message, slog.String("err", fmt.Sprintf("%v", value)))
}

func (l *Log) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *Log) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type loggerImpl struct {
	*slog.Logger
}

var Log Logger

func Init() {
	Log = &loggerImpl{slog.New(slog.NewJSONHandler(os.Stdout, nil))}
}

func (l *loggerImpl) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

func (l *loggerImpl) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

func (l *loggerImpl) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

func (l *loggerImpl) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

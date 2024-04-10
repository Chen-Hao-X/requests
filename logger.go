package requests

import (
	"fmt"
	"io"

	"golang.org/x/exp/slog"
)

type Level int

// The levels of logs.
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger is a logger interface that output logs with a format.
type Logger interface {
	Debugf(format string, v ...any)
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
	SetLevel(level Level)
}

type DefaultLogger struct {
	logger *slog.Logger
	level  *slog.LevelVar
}

// Debugf logs a message at the Debug level.
func (l *DefaultLogger) Debugf(format string, v ...any) {
	l.logger.Debug(fmt.Sprintf(format, v...))
}

// Infof logs a message at the Info level.
func (l *DefaultLogger) Infof(format string, v ...any) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

// Warnf logs a message at the Warn level.
func (l *DefaultLogger) Warnf(format string, v ...any) {
	l.logger.Warn(fmt.Sprintf(format, v...))
}

// Errorf logs a message at the Error level.
func (l *DefaultLogger) Errorf(format string, v ...any) {
	l.logger.Error(fmt.Sprintf(format, v...))
}

// SetLevel sets the log level of the logger.
func (l *DefaultLogger) SetLevel(level Level) {
	switch level {
	case LevelDebug:
		l.level.Set(slog.LevelDebug)
	case LevelInfo:
		l.level.Set(slog.LevelInfo)
	case LevelWarn:
		l.level.Set(slog.LevelWarn)
	case LevelError:
		l.level.Set(slog.LevelError)
	}
}

func NewDefaultLogger(output io.Writer, level Level) Logger {
	levelVar := &slog.LevelVar{}

	// Initialize text handler at the desired log level if possible
	textHandler := slog.NewTextHandler(output, &slog.HandlerOptions{
		Level: levelVar,
	})

	// Create and return a new `DefaultLogger`
	logger := &DefaultLogger{
		logger: slog.New(textHandler),
		level:  levelVar,
	}

	logger.SetLevel(level)
	return logger
}

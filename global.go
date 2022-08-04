package golog

import (
	"sync"

	"github.com/millken/golog/internal/config"
	"github.com/millken/golog/internal/log"
)

//nolint:gochecknoglobals
var (
	loggerProviderInstance log.Logger
	loggerProviderOnce     sync.Once
)

func F(k string, v interface{}) log.Field {
	return log.Field{Key: k, Val: v}
}

// Panicf logs a message using Panic level and panics.
func Panicf(format string, args ...interface{}) {
	loggerProvider().Panicf(format, args...)
}

// Fatalf logs a message using Fatal level and exits with status 1.
func Fatalf(format string, args ...interface{}) {
	loggerProvider().Fatalf(format, args...)
}

// Errorf logs a message using Error level.
func Errorf(format string, args ...interface{}) {
	loggerProvider().Errorf(format, args...)
}

// Warnf logs a message using Warn level.
func Warnf(format string, args ...interface{}) {
	loggerProvider().Warnf(format, args...)
}

// Infof logs a message using Info level.
func Infof(format string, args ...interface{}) {
	loggerProvider().Infof(format, args...)
}

// Debugf logs a message using Debug level.
func Debugf(format string, args ...interface{}) {
	loggerProvider().Debugf(format, args...)
}

// WithField returns a logger configured with the key-value pair.
func WithField(k string, v interface{}) log.Logger {
	return WithFields(log.Field{Key: k, Val: v})
}

// WithFields returns a logger configured with the key-value pairs.
func WithFields(fields ...log.Field) log.Logger {
	return loggerProvider().WithFields(fields...)
}

// LoadConfig - Load configuration from file
func LoadConfig(path string) error {
	return config.Load(path)
}

// Panic logs a message using Panic level and panics.
func Panic(msg string) {
	loggerProvider().Panic(msg)
}

// Fatal logs a message using Fatal level and exits with status 1.
func Fatal(msg string) {
	loggerProvider().Fatal(msg)
}

// Error logs a message using Error level.
func Error(msg string) {
	loggerProvider().Error(msg)
}

// Warn logs a message using Warn level.
func Warn(msg string) {
	loggerProvider().Warn(msg)
}

// Info logs a message using Info level.
func Info(msg string) {
	loggerProvider().Infof(msg)
}

// Debug logs a message using Debug level.
func Debug(msg string) {
	loggerProvider().Debugf(msg)
}

func loggerProvider() log.Logger {
	loggerProviderOnce.Do(func() {
		loggerProviderInstance = New("_global")
	})

	return loggerProviderInstance
}

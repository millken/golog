package golog

import (
	"sync"
	"sync/atomic"
)

type loggerProviderFactory func() Logger

//nolint:gochecknoglobals
var (
	loggerProviderFactoryFn atomic.Value
)

func init() {
	loggerProviderFactoryFn.Store(newLoggerProviderFactory())
}

func newLoggerProviderFactory() loggerProviderFactory {
	return sync.OnceValue(func() Logger {
		return New(defaultModule).CallerSkip(1)
	})
}

// field is a shortcut to create Field.
func field(k string, v any) Field {
	return Field{Key: k, Val: v}
}

// Panicf logs a message using Panic level and panics.
func Panicf(format string, args ...any) {
	loggerProvider().Panicf(format, args...)
}

// Fatalf logs a message using Fatal level and exits with status 1.
func Fatalf(format string, args ...any) {
	loggerProvider().Fatalf(format, args...)
}

// Errorf logs a message using Error level.
func Errorf(format string, args ...any) {
	loggerProvider().Errorf(format, args...)
}

// Warnf logs a message using Warn level.
func Warnf(format string, args ...any) {
	loggerProvider().Warnf(format, args...)
}

// Infof logs a message using Info level.
func Infof(format string, args ...any) {
	loggerProvider().Infof(format, args...)
}

// Debugf logs a message using Debug level.
func Debugf(format string, args ...any) {
	loggerProvider().Debugf(format, args...)
}

// WithValues returns a logger configured with the key-value pairs.
func WithValues(keysAndVals ...any) Logger {
	return loggerProvider().WithValues(keysAndVals...)
}

// Panic logs a message using Panic level and panics.
func Panic(msg string, keysAndVals ...any) {
	loggerProvider().Panic(msg, keysAndVals...)
}

// Fatal logs a message using Fatal level and exits with status 1.
func Fatal(msg string, keysAndVals ...any) {
	loggerProvider().Fatal(msg, keysAndVals...)
}

// Error logs a message using Error level.
func Error(msg string, keysAndVals ...any) {
	loggerProvider().Error(msg, keysAndVals...)
}

// Warn logs a message using Warn level.
func Warn(msg string, keysAndVals ...any) {
	loggerProvider().Warn(msg, keysAndVals...)
}

// Info logs a message using Info level.
func Info(msg string, keysAndVals ...any) {
	loggerProvider().Info(msg, keysAndVals...)
}

// Debug logs a message using Debug level.
func Debug(msg string, keysAndVals ...any) {
	loggerProvider().Debug(msg, keysAndVals...)
}

func loggerProvider() Logger {
	f := loggerProviderFactoryFn.Load().(loggerProviderFactory)
	return f()
}

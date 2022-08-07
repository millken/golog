package golog

import (
	"sync"
)

//nolint:gochecknoglobals
var (
	loggerProviderInstance Logger
	loggerProviderOnce     sync.Once
)

// field is a shortcut to create Field.
func field(k string, v interface{}) Field {
	return Field{Key: k, Val: v}
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

// WithValues returns a logger configured with the key-value pairs.
func WithValues(keysAndVals ...interface{}) Logger {
	return loggerProvider().WithValues(keysAndVals)
}

// Panic logs a message using Panic level and panics.
func Panic(msg string, keysAndVals ...interface{}) {
	loggerProvider().Panic(msg, keysAndVals...)
}

// Fatal logs a message using Fatal level and exits with status 1.
func Fatal(msg string, keysAndVals ...interface{}) {
	loggerProvider().Fatal(msg, keysAndVals...)
}

// Error logs a message using Error level.
func Error(msg string, keysAndVals ...interface{}) {
	loggerProvider().Error(msg, keysAndVals...)
}

// Warn logs a message using Warn level.
func Warn(msg string, keysAndVals ...interface{}) {
	loggerProvider().Warn(msg, keysAndVals...)
}

// Info logs a message using Info level.
func Info(msg string, keysAndVals ...interface{}) {
	loggerProvider().Info(msg, keysAndVals...)
}

// Debug logs a message using Debug level.
func Debug(msg string, keysAndVals ...interface{}) {
	loggerProvider().Debug(msg, keysAndVals...)
}

func loggerProvider() Logger {
	loggerProviderOnce.Do(func() {
		loggerProviderInstance = New("").CallerSkip(1)
	})

	return loggerProviderInstance
}

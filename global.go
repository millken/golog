package golog

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	// TimestampFieldName is the field name used for the timestamp field.
	TimestampFieldName = "time"

	// LevelFieldName is the field name used for the level field.
	LevelFieldName = "level"

	// LevelFieldMarshalFunc allows customization of global level field marshaling
	LevelFieldMarshalFunc = func(l Level) string {
		return *l.String()
	}

	// MessageFieldName is the field name used for the message field.
	MessageFieldName = "message"

	// ErrorFieldName is the field name used for error fields.
	ErrorFieldName = "error"

	// CallerFieldName is the field name used for caller field.
	CallerFieldName = "caller"

	// ErrorStackFieldName is the field name used for error stacks.
	ErrorStackFieldName = "stack"

	// ErrorStackMarshaler extract the stack from err if any.
	ErrorStackMarshaler func(err error) interface{}

	// ErrorMarshalFunc allows customization of global error marshaling
	ErrorMarshalFunc = func(err error) interface{} {
		return err
	}

	// TimeFieldFormat defines the time format of the Time field type. If set to
	// TimeFormatUnix, TimeFormatUnixMs or TimeFormatUnixMicro, the time is formatted as an UNIX
	// timestamp as integer.
	TimeFieldFormat = time.RFC3339

	// TimestampFunc defines the function called to generate a timestamp.
	TimestampFunc = time.Now

	// DurationFieldUnit defines the unit for time.Duration type fields added
	// using the Dur method.
	DurationFieldUnit = time.Millisecond

	// DurationFieldInteger renders Dur fields as integer instead of float if
	// set to true.
	DurationFieldInteger = false

	// ErrorHandler is called whenever zerolog fails to write an event on its
	// output. If not set, an error is printed on the stderr. This handler must
	// be thread safe and non-blocking.
	ErrorHandler func(err error)
)

// Level defines log levels.
type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// Disabled disables the logger.
	Disabled

	// DefaultLevel defines the default log level.
	DefaultLevel = InfoLevel
)

var (
	levelMessages = []string{
		DebugLevel: "DEBUG",
		InfoLevel:  "INFO",
		WarnLevel:  "WARN",
		ErrorLevel: "ERROR",
		FatalLevel: "FATAL",
		PanicLevel: "PANIC",
	}
)

//using pointer can reduce allocs
func (l Level) String() *string {
	return &levelMessages[l]
}

// ParseLevel takes a string level and returns the Logrus log level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "panic":
		return PanicLevel, nil
	case "fatal":
		return FatalLevel, nil
	case "error":
		return ErrorLevel, nil
	case "warn", "warning":
		return WarnLevel, nil
	case "info":
		return InfoLevel, nil
	case "debug":
		return DebugLevel, nil
	}

	var l Level
	return l, fmt.Errorf("not a valid Level: %q", lvl)
}

// Levels returns the available logging levels.
type Levels []Level

// Contains returns true if the slice contains the level.
func (l Levels) Contains(level Level) bool {
	for _, lv := range l {
		if lv == level {
			return true
		}
	}
	return false
}

var (
	_globalMu     sync.RWMutex
	_globalLogger = NewLogger()
)

// ReplaceGlobals replaces the global Logger, and returns a
// function to restore the original values. It's safe for concurrent use.
func ReplaceGlobals(logger *Logger) func() {
	_globalMu.Lock()
	prev := _globalLogger
	_globalLogger = logger
	_globalMu.Unlock()
	return func() { ReplaceGlobals(prev) }
}

// Panicf logs a message using Panic level and panics.
func Panicf(format string, args ...interface{}) {
	safeLogger().Panicf(format, args...)
}

// Fatalf logs a message using Fatal level and exits with status 1.
func Fatalf(format string, args ...interface{}) {
	safeLogger().Fatalf(format, args...)
}

// Errorf logs a message using Error level.
func Errorf(format string, args ...interface{}) {
	safeLogger().Errorf(format, args...)
}

// Warnf logs a message using Warn level.
func Warnf(format string, args ...interface{}) {
	safeLogger().Warnf(format, args...)
}

// Infof logs a message using Info level.
func Infof(format string, args ...interface{}) {
	safeLogger().Infof(format, args...)
}

// Debugf logs a message using Debug level.
func Debugf(format string, args ...interface{}) {
	safeLogger().Debugf(format, args...)
}

// Panic logs a message using Panic level and panics.
func Panic(msg string) {
	safeLogger().Panic(msg)
}

// Fatal logs a message using Fatal level and exits with status 1.
func Fatal(msg string) {
	safeLogger().Fatal(msg)
}

// Error logs a message using Error level.
func Error(msg string) {
	safeLogger().Error(msg)
}

// Warn logs a message using Warn level.
func Warn(msg string) {
	safeLogger().Warn(msg)
}

// Info logs a message using Info level.
func Info(msg string) {
	safeLogger().Infof(msg)
}

// Debug logs a message using Debug level.
func Debug(msg string) {
	safeLogger().Debugf(msg)
}

// WithField returns a logger configured with the key-value pair.
func WithField(k string, v interface{}) *Logger {
	return WithFields(Field{Key: k, Val: v})
}

// WithFields returns a logger configured with the key-value pairs.
func WithFields(fields ...Field) *Logger {
	l := safeLogger()
	l.CallerSkip = l.CallerSkip - 1
	return l.WithFields(fields...)
}

// safeLogger returns the global Logger, which can be reconfigured with ReplaceGlobals.
// It's safe for concurrent use.
func safeLogger() *Logger {
	_globalMu.RLock()
	l := _globalLogger
	_globalMu.RUnlock()
	l.CallerSkip = 2
	return l
}

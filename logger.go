// Package golog provides a high-performance structured logging library
// with support for multiple output formats, log rotation, and module-based
// configuration.
//
// Basic usage:
//
//	golog.Info("Hello, world!", "key", "value")
//	golog.Errorf("Error occurred: %v", err)
//
// Creating a module-specific logger:
//
//	logger := golog.New("my/module")
//	logger.Info("Module message", "user", "john")
//
// Configuration:
//
//	golog.SetLevel(golog.DEBUG)
//	golog.SetEncoding(golog.JSONEncoding)
package golog

import (
	"time"
)

const (
	// TimestampFieldName is the field name used for the timestamp field.
	TimestampFieldName = "time"

	// LevelFieldName is the field name used for the level field.
	LevelFieldName = "level"

	// ModuleFieldName is the field name used for the module field.
	ModuleFieldName = "module"
	// MessageFieldName is the field name used for the message field.
	MessageFieldName = "message"

	// ErrorFieldName is the field name used for error fields.
	ErrorFieldName = "error"

	// CallerFieldName is the field name used for caller field.
	CallerFieldName = "caller"

	// ErrorStackFieldName is the field name used for error stacks.
	ErrorStackFieldName = "stack"

	// TimeFieldFormat defines the time format of the Time field type. If set to
	// TimeFormatUnix, TimeFormatUnixMs or TimeFormatUnixMicro, the time is formatted as an UNIX
	// timestamp as integer.
	TimeFieldFormat = time.RFC3339
)

// Field is a key/value pair.
type Field struct {
	Key string
	Val any
}

// Logger represents a general-purpose logger.
type Logger interface {
	WithValues(keysAndVals ...any) Logger
	Panicf(msg string, args ...any)
	Fatalf(msg string, args ...any)
	Errorf(msg string, args ...any)
	Warnf(msg string, args ...any)
	Infof(msg string, args ...any)
	Debugf(msg string, args ...any)
	Panic(msg string, keysAndVals ...any)
	Fatal(msg string, keysAndVals ...any)
	Error(msg string, keysAndVals ...any)
	Warn(msg string, keysAndVals ...any)
	Info(msg string, keysAndVals ...any)
	Debug(msg string, keysAndVals ...any)
}

// Encoder is an interface for encoding log entry.
type Encoder interface {
	Encode(*Entry) ([]byte, error)
}

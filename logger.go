package golog

import (
	"time"
)

const (
	// TimestampFieldName is the field name used for the timestamp field.
	TimestampFieldName = "time"

	// LevelFieldName is the field name used for the level field.
	LevelFieldName = "level"

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
	Val interface{}
}

// Fields is a map of fields.
type Fields map[string]interface{}

// Logger represents a general-purpose logger.
type Logger interface {
	WithValues(keysAndVals ...interface{}) Logger
	Panicf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Panic(msg string, keysAndVals ...interface{})
	Fatal(msg string, keysAndVals ...interface{})
	Error(msg string, keysAndVals ...interface{})
	Warn(msg string, keysAndVals ...interface{})
	Info(msg string, keysAndVals ...interface{})
	Debug(msg string, keysAndVals ...interface{})
}

// Encoder is a interface for encoding log entry.
type Encoder interface {
	Encode(*Entry) ([]byte, error)
}

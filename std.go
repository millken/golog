package golog

import (
	"io"
	"os"
	"time"
)

const defaultTextSeparator = " - "

const (
	defaultJSONFieldKeyDatetime  = "datetime"
	defaultJSONFieldKeyTimestamp = "timestamp"
	defaultJSONFieldKeyLevel     = "level"
	defaultJSONFieldKeyFile      = "file"
	defaultJSONFieldKeyStack     = "stack"
	defaultJSONFieldKeyMessage   = "message"
)

const defaultDatetimeLayout = time.RFC3339

const defaultTimestampFormat = TimestampFormatSeconds

// Field is a key/value pair.
type Field struct {
	Key string
	Val interface{}
}

var std = newStd()

func newStd() *Logger {
	defaultConfigValue.Store(config)
	defaultHandlerValue.Store(newDefaultHandler(os.Stdout))
	l := New()
	l.setCalldepth(calldepthStd)

	return l
}

// Panicf logs a message using Panic level and panics.
func Panicf(format string, args ...interface{}) {
	std.Panicf(format, args...)
}

// Fatalf logs a message using Fatal level and exits with status 1.
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}

// Errorf logs a message using Error level.
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

// Warnf logs a message using Warn level.
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}

// Infof logs a message using Info level.
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}

// Debugf logs a message using Debug level.
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}

// WithValues returns a logger configured with the key-value pairs.
func WithValues(keysAndVals ...interface{}) *Logger {
	return std.WithValues(keysAndVals...)
}

// Panic logs a message using Panic level and panics.
func Panic(msg string, keysAndVals ...interface{}) {
	std.Panic(msg, keysAndVals...)
}

// Fatal logs a message using Fatal level and exits with status 1.
func Fatal(msg string, keysAndVals ...interface{}) {
	std.Fatal(msg, keysAndVals...)
}

// Error logs a message using Error level.
func Error(msg string, keysAndVals ...interface{}) {
	std.Error(msg, keysAndVals...)
}

// Warn logs a message using Warn level.
func Warn(msg string, keysAndVals ...interface{}) {
	std.Warn(msg, keysAndVals...)
}

// Info logs a message using Info level.
func Info(msg string, keysAndVals ...interface{}) {
	std.Info(msg, keysAndVals...)
}

// Debug logs a message using Debug level.
func Debug(msg string, keysAndVals ...interface{}) {
	std.Debug(msg, keysAndVals...)
}

// AddHandler adds a handler to the default logger.
func AddHandler(h Handler) {
	std.AddHandler(h)
}

// SetEncoder sets the encoder to the standard logger.
func SetEncoder(enc Encoder) {
	switch t := enc.(type) {
	case *EncoderJSON:
		defaultEncoderType.Store(true)
		defaultEncoderJson.Store(t)
	case *EncoderText:
		defaultEncoderType.Store(false)
		defaultEncoderText.Store(t)
	}
}

// SetWriter sets the writer to the standard logger.
func SetWriter(w io.Writer) {
	defaultHandler().SetWriter(w)
}

// SetLevel sets the level to the standard logger.
func SetLevel(lvl Level) {
	cfg := defaultConfig()
	cfg.Level = lvl
	defaultConfigValue.Store(cfg)
}

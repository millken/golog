package golog

import (
	"fmt"
	"os"
	"time"

	"github.com/millken/golog/internal/log"
	"github.com/millken/golog/internal/meta"
)

var (
	_ log.Logger = (*Log)(nil)
)

// Log is an implementation of Logger interface.
// It encapsulates default or custom logger to provide module and level based logging.
type Log struct {
	module string
	fields []log.Field
}

// New creates and returns a Logger implementation based on given module name.
// note: the underlying logger instance is lazy initialized on first use.
// To use your own logger implementation provide logger provider in 'Initialize()' before logging any line.
// If 'Initialize()' is not called before logging any line then default logging implementation will be used.
func NewLog(module string) *Log {
	return &Log{module: module}
}

// Fatalf calls underlying logger.Fatal.
func (l *Log) Fatalf(format string, args ...interface{}) {
	if !meta.IsEnabledFor(l.module, log.FATAL) {
		return
	}

	l.logf(log.FATAL, format, args...)
	os.Exit(1)
}

// Panicf calls underlying logger.Panic.
func (l *Log) Panicf(format string, args ...interface{}) {
	if !meta.IsEnabledFor(l.module, log.PANIC) {
		return
	}

	l.logf(log.PANIC, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Debugf calls error log function if DEBUG level enabled.
func (l *Log) Debugf(format string, args ...interface{}) {
	if !meta.IsEnabledFor(l.module, log.DEBUG) {
		return
	}

	l.logf(log.DEBUG, format, args...)
}

// Infof calls error log function if INFO level enabled.
func (l *Log) Infof(format string, args ...interface{}) {
	if !meta.IsEnabledFor(l.module, log.INFO) {
		return
	}

	l.logf(log.INFO, format, args...)
}

// Warnf calls error log function if WARNING level enabled.
func (l *Log) Warnf(format string, args ...interface{}) {
	if !meta.IsEnabledFor(l.module, log.WARNING) {
		return
	}

	l.logf(log.WARNING, format, args...)
}

// Errorf calls error log function if ERROR level enabled.
func (l *Log) Errorf(format string, args ...interface{}) {
	if !meta.IsEnabledFor(l.module, log.ERROR) {
		return
	}

	l.logf(log.ERROR, format, args...)
}

// Fatal calls underlying logger.Fatal.
func (l *Log) Fatal(msg string) {
	if !meta.IsEnabledFor(l.module, log.FATAL) {
		return
	}

	l.logf(log.FATAL, msg)
	os.Exit(1)
}

// Panic calls underlying logger.Panic.
func (l *Log) Panic(msg string) {
	if !meta.IsEnabledFor(l.module, log.PANIC) {
		return
	}

	l.logf(log.PANIC, msg)
	panic(msg)
}

// Debug calls error log function if DEBUG level enabled.
func (l *Log) Debug(msg string) {
	if !meta.IsEnabledFor(l.module, log.DEBUG) {
		return
	}

	l.logf(log.DEBUG, msg)
}

// Info calls error log function if INFO level enabled.
func (l *Log) Info(msg string) {
	if !meta.IsEnabledFor(l.module, log.INFO) {
		return
	}

	l.logf(log.INFO, msg)
}

// Warn calls error log function if WARNING level enabled.
func (l *Log) Warn(msg string) {
	if !meta.IsEnabledFor(l.module, log.WARNING) {
		return
	}

	l.logf(log.WARNING, msg)
}

// Error calls error log function if ERROR level enabled.
func (l *Log) Error(msg string) {
	if !meta.IsEnabledFor(l.module, log.ERROR) {
		return
	}

	l.logf(log.ERROR, msg)
}

// WithField returns a logger configured with the key-value pair.
func (l *Log) WithField(k string, v interface{}) log.Logger {
	return &Log{fields: append(l.fields, log.Field{Key: k, Val: v}), module: l.module}
}

// WithFields returns a logger configured with the key-value pairs.
func (l *Log) WithFields(fields ...log.Field) log.Logger {
	return &Log{fields: append(l.fields, fields...), module: l.module}
}

func (l *Log) logf(level log.Level, format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.output(level, msg, l.fields...)
}

func (l *Log) output(level log.Level, msg string, fields ...log.Field) {
	entry := log.AcquireEntry()
	entry.Module = l.module
	handler := meta.GetHandler(l.module)
	copy(entry.Fields[0:len(fields)], fields)
	entry.SetFieldsLen(len(fields))

	entry.Message = msg
	entry.Level = level
	entry.Timestamp = time.Now()
	if err := handler.Handle(entry); err != nil {
		fmt.Fprintf(os.Stderr, "golog: failed to handle entry: %v", err)
	}
	entry.Reset()
	log.ReleaseEntry(entry)
}

// SetLevel - setting log level for given module
//  Parameters:
//  module is module name
//  level is logging level
//
// If not set default logging level is info.
func SetLevel(module string, level log.Level) {
	meta.SetLevel(module, level)
}

// GetLevel - getting log level for given module
//  Parameters:
//  module is module name
//
//  Returns:
//  logging level
//
// If not set default logging level is info.
func GetLevel(module string) log.Level {
	return meta.GetLevel(module)
}

// IsEnabledFor - Check if given log level is enabled for given module
//  Parameters:
//  module is module name
//  level is logging level
//
//  Returns:
//  is logging enabled for this module and level
//
// If not set default logging level is info.
func IsEnabledFor(module string, level log.Level) bool {
	return meta.IsEnabledFor(module, level)
}

// ParseLevel returns the log level from a string representation.
//  Parameters:
//  level is logging level in string representation
//
//  Returns:
//  logging level
func ParseLevel1(level string) (log.Level, error) {
	l, err := meta.ParseLevel(level)

	return l, err
}

// ShowCallerInfo - Show caller info in log lines for given log level and module
//  Parameters:
//  module is module name
//  level is logging level
//
// note: based on implementation of custom logger, callerinfo info may not be available for custom logging provider
func ShowCallerInfo(module string, level log.Level) {
	meta.ShowCallerInfo(module, level)
}

// HideCallerInfo - Do not show caller info in log lines for given log level and module
//  Parameters:
//  module is module name
//  level is logging level
//
// note: based on implementation of custom logger, callerinfo info may not be available for custom logging provider
func HideCallerInfo(module string, level log.Level) {
	meta.HideCallerInfo(module, level)
}

// IsCallerInfoEnabled - returns if caller info enabled for given log level and module
//  Parameters:
//  module is module name
//  level is logging level
//
//  Returns:
//  is caller info enabled for this module and level
//
// note: based on implementation of custom logger, callerinfo info may not be available for custom logging provider
func IsCallerInfoEnabled(module string, level log.Level) bool {
	return meta.IsCallerInfoEnabled(module, level)
}

package log

import (
	"io"
	"os"
	"sync/atomic"

	"github.com/millken/golog"
)

var (
	handler      golog.Handler
	logger       *golog.Logger
	formatter    golog.Formatter
	output       io.Writer
	level        golog.Level
	renew        int32
	color        bool
	enableCaller bool
	enableStack  bool
)

type log struct {
	*golog.Logger
}

func init() {
	SetFormatter(golog.NewTextFormatter())
	SetOutput(os.Stdout)
	SetLevel(golog.InfoLevel)
}

func s() *golog.Logger {
	if atomic.LoadInt32(&renew) == 1 {
		handler = golog.NewLoggerHandler(output)
		handler.SetFormatter(formatter)
		handler.SetLevel(level)
		logger = golog.NewLogger()
		logger.AddHandler(handler)
		atomic.StoreInt32(&renew, 0)
	}
	logger.CallerSkip = 2
	return logger
}

// SetFormatter sets the formatter for the logger.
func SetFormatter(f golog.Formatter) {
	formatter = f
	atomic.StoreInt32(&renew, 1)
}

// SetOutput sets the output for the logger.
func SetOutput(w io.Writer) {
	output = w
	atomic.StoreInt32(&renew, 1)
}

// SetLevel sets the level for the logger.
func SetLevel(l golog.Level) {
	level = l
	atomic.StoreInt32(&renew, 1)
}

// WithField returns a new entry with the field added to it.
func WithField(k string, v interface{}) *golog.Logger {
	l := s()
	l.CallerSkip--
	return l.WithField(k, v)
}

// WithFields returns a new entry with the fields added to it.
func WithFields(fields ...golog.Field) *golog.Logger {
	l := s()
	l.CallerSkip--
	return l.WithFields(fields...)
}

// Debugf logs a message at debug level.
func Debugf(format string, args ...interface{}) {
	s().Debugf(format, args...)
}

// Infof logs a message at info level.
func Infof(format string, args ...interface{}) {
	s().Infof(format, args...)
}

// Warnf logs a message at warn level.
func Warnf(format string, args ...interface{}) {
	s().Warnf(format, args...)
}

// Errorf logs a message at error level.
func Errorf(format string, args ...interface{}) {
	s().Errorf(format, args...)
}

// Fatalf logs a message at fatal level.
func Fatalf(format string, args ...interface{}) {
	s().Fatalf(format, args...)
}

// Panicf logs a message at panic level.
func Panicf(format string, args ...interface{}) {
	s().Panicf(format, args...)
}

// Debug logs a message at debug level.
func Debug(msg string) {
	s().Debug(msg)
}

// Info logs a message at info level.
func Info(msg string) {
	s().Info(msg)
}

// Warn logs a message at warn level.
func Warn(msg string) {
	s().Warn(msg)
}

// Error logs a message at error level.
func Error(msg string) {
	s().Error(msg)
}

// Fatal logs a message at fatal level.
func Fatal(msg string) {
	s().Fatal(msg)
}

// Panic logs a message at panic level.
func Panic(msg string) {
	s().Panic(msg)
}

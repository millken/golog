package golog

import (
	"os"
	"runtime"
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

	// CallerSkipFrameCount is the number of stack frames to skip to find the caller.
	CallerSkipFrameCount = 5

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
	// NoLevel defines an absent log level.
	NoLevel Level = iota
	// DebugLevel defines debug log level.
	DebugLevel
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

	DefaultLevel = InfoLevel
)

var (
	levelMessages = []string{
		DebugLevel: "debug",
		InfoLevel:  "info",
		WarnLevel:  "warn",
		ErrorLevel: "error",
		FatalLevel: "fatal",
		PanicLevel: "panic",
	}

	_cwd, _ = os.Getwd()
)

//using pointer can reduce allocs
func (l Level) String() *string {
	return &levelMessages[l]
}

type Levels []Level

func (l Levels) Contains(level Level) bool {
	for _, lv := range l {
		if lv == level {
			return true
		}
	}
	return false
}

func caller(skip int) (string, int) {
	_, file, line, _ := runtime.Caller(skip)
	return file, line
}

var (
	_globalMu sync.RWMutex
	_globalL  = NewLogger()
)

// ReplaceGlobals replaces the global Logger, and returns a
// function to restore the original values. It's safe for concurrent use.
func ReplaceGlobals(logger *Logger) func() {
	_globalMu.Lock()
	prev := _globalL
	_globalL = logger
	_globalMu.Unlock()
	return func() { ReplaceGlobals(prev) }
}

func Fatalf(format string, args ...interface{}) {
	safeLogger().Fatalf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	safeLogger().Errorf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	safeLogger().Warnf(format, args...)
}

func Infof(format string, args ...interface{}) {
	safeLogger().Infof(format, args...)
}

func Debugf(format string, args ...interface{}) {
	safeLogger().Debugf(format, args...)
}

func WithField(k string, v interface{}) *Logger {
	return safeLogger().WithFields(field{k, v})
}

func WithFields(fields ...field) *Logger {
	return safeLogger().WithFields(fields...)
}

func AddHandler(handler Handler) {
	safeLogger().AddHandler(handler)
}

// safeLogger returns the global Logger, which can be reconfigured with ReplaceGlobals.
// It's safe for concurrent use.
func safeLogger() *Logger {
	_globalMu.RLock()
	l := _globalL
	_globalMu.RUnlock()
	return l
}

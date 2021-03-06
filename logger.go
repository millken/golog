package golog

import (
	"fmt"
	"os"
	"time"
)

// Logger is a simple logger.
type Logger struct {
	callerSkip int
	handlers   []Handler
	fields     []Field
}

// NewLogger creates a new Logger.
func NewLogger() *Logger {
	log := &Logger{
		handlers: make([]Handler, 0),
		fields:   make([]Field, 0, 512),
	}
	return log
}

func (l *Logger) output(level Level, msg string, fields ...Field) {
	if len(l.handlers) == 0 {
		return
	}
	for _, handler := range l.handlers {
		if len(handler.Levels()) > 0 {
			if !handler.Levels().Contains(level) {
				continue
			}
		} else if handler.Level() > level {
			continue
		}
		entry := acquireEntry()
		if !handler.DisableLogFields() {
			copy(entry.Fields[0:len(fields)], fields)
			entry.fieldsLen = len(fields)
		}

		entry.Message = msg
		entry.Level = level
		entry.Timestamp = time.Now()
		entry.callerSkip = l.callerSkip + 3
		entry.Reset()

		formatter := handler.Formatter()
		if formatter != nil {
			err := formatter.Format(entry)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		if err := handler.Handle(entry); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		releaseEntry(entry)
	}
}

// AddHandler adds a handler.
func (l *Logger) AddHandler(handler Handler) {
	l.callerSkip++
	l.handlers = append(l.handlers, handler)
}

// WithField returns a new logger with the field added.
func (l *Logger) WithField(k string, v interface{}) *Logger {
	l.callerSkip++
	return l.WithFields(Field{k, v})
}

// WithFields returns a new logger with the fields added.
func (l *Logger) WithFields(fields ...Field) *Logger {
	return &Logger{
		callerSkip: 1,
		handlers:   l.handlers,
		fields:     append(l.fields, fields...),
	}
}

func (l *Logger) logf(level Level, format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.output(level, msg, l.fields...)
}

// Debugf logs a message at debug level.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(DebugLevel, format, args...)
}

// Infof logs a message at info level.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(InfoLevel, format, args...)
}

// Warnf logs a message at warn level.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logf(WarnLevel, format, args...)
}

// Errorf logs a message at error level.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(ErrorLevel, format, args...)
}

// Fatalf logs a message using Fatal level and exits with status 1.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(FatalLevel, format, args...)
	os.Exit(1)
}

// Panicf logs a message using Panic level and panics.
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.logf(PanicLevel, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Reset resets the logger.
func (l *Logger) Reset() {
	l.handlers = l.handlers[:0]
	l.fields = l.fields[:0]
}

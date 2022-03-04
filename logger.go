package golog

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	handlers []Handler
	fields   []field
}

func NewLogger() *Logger {
	log := &Logger{}
	return log
}

func (l *Logger) output(level Level, msg string, fields ...field) error {
	var err error
	for _, handler := range l.handlers {
		if handler.GetLevel() > level {
			continue
		}
		entry := acquireEntry()
		copy(entry.Fields[0:len(fields)], fields)
		entry.fieldsLen = len(fields)

		entry.Message = msg
		entry.Level = level
		entry.Timestamp = time.Now()
		entry.Reset()

		formatter := handler.GetFormatter()
		if formatter != nil {
			err = formatter.Format(entry)
			if err != nil {
				fmt.Println(err)
			}
		}
		handler.Handle(entry)
		releaseEntry(entry)
	}
	return err
}

func (l *Logger) AddHandler(handler Handler) {
	l.handlers = append(l.handlers, handler)
}

func (l *Logger) Debug(msg string, fields ...field) {
	l.output(DebugLevel, msg, fields...)
}

func (l *Logger) Info(msg string, fields ...field) {
	l.output(InfoLevel, msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...field) {
	l.output(WarnLevel, msg, fields...)
}

func (l *Logger) Error(msg string, fields ...field) {
	l.output(ErrorLevel, msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...field) {
	l.output(FatalLevel, msg, fields...)
	os.Exit(1)
}

func (l *Logger) WithField(k string, v interface{}) *Logger {
	return l.WithFields(field{k, v})
}

func (l *Logger) WithFields(fields ...field) *Logger {
	return &Logger{
		handlers: l.handlers,
		fields:   append(l.fields, fields...),
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

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(DebugLevel, format, args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(InfoLevel, format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logf(WarnLevel, format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(ErrorLevel, format, args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logf(FatalLevel, format, args...)
	os.Exit(1)
}

func (l *Logger) Reset() {
	l.handlers = l.handlers[:0]
	l.fields = l.fields[:0]
}

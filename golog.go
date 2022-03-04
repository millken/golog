package golog

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	*logger
}

type logger struct {
	handlers []Handler
	fields   []field
}

func NewLogger() *Logger {
	log := &Logger{
		newLogger(),
	}
	return log
}

func newLogger() *logger {
	return &logger{}
}

func (l *logger) output(level Level, msg string, fields ...field) error {
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

func (l *logger) AddHandler(handler Handler) {
	l.handlers = append(l.handlers, handler)
}

func (l *logger) Debug(msg string, fields ...field) {
	l.output(DebugLevel, msg, fields...)
}

func (l *logger) Info(msg string, fields ...field) {
	l.output(InfoLevel, msg, fields...)
}

func (l *logger) Warn(msg string, fields ...field) {
	l.output(WarnLevel, msg, fields...)
}

func (l *logger) Error(msg string, fields ...field) {
	l.output(ErrorLevel, msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...field) {
	l.output(FatalLevel, msg, fields...)
	os.Exit(1)
}

func (l *logger) WithField(k string, v interface{}) *logger {
	return l.WithFields(field{k, v})
}

func (l *logger) WithFields(fields ...field) *logger {
	return &logger{
		handlers: l.handlers,
		fields:   append(l.fields, fields...),
	}
}

func (l *logger) logf(level Level, format string, args ...interface{}) {
	l.output(level, fmt.Sprintf(format, args...), l.fields...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.logf(DebugLevel, format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.logf(InfoLevel, format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.logf(WarnLevel, format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.logf(ErrorLevel, format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.logf(FatalLevel, format, args...)
	os.Exit(1)
}

func (l *logger) Reset() {
	l.handlers = l.handlers[:0]
	l.fields = l.fields[:0]
}

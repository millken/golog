package golog

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/millken/golog/pool"
)

// Event pool
var loggerPool = pool.NewReferenceCountedPool(
	func(counter pool.ReferenceCounter) pool.ReferenceCountable {
		br := new(logger)
		br.ReferenceCounter = counter
		return br
	}, resetLogger)

// Method to get new Event
func acquireLogger() *logger {
	return loggerPool.Get().(*logger)
}

// Method to reset Event
// Used by reference countable pool
func resetLogger(i interface{}) error {
	obj, ok := i.(*logger)
	if !ok {
		errors.Errorf("illegal object sent to ResetEvent %v", i)
	}
	obj.Reset()
	return nil
}

//easyjson:json
type Fields map[string]interface{}

type Logger struct {
	*logger
}

type logger struct {
	pool.ReferenceCounter
	handlers []Handler
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
		defer releaseEntry(entry)
		entry.Fields = fields

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
}

func (l *logger) Reset() {
	l.handlers = l.handlers[:0]
}

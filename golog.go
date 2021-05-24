package golog

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/millken/golog/pool"
	"github.com/valyala/bytebufferpool"
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
	fields   fields
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

func (l *logger) output(level Level, msg ...interface{}) error {

	buff := bytebufferpool.Get()
	defer bytebufferpool.Put(buff)
	buff.Reset()
	fmt.Fprint(buff, msg...)

	var err error
	for _, handler := range l.handlers {
		if handler.GetLevel() > level {
			continue
		}
		entry := acquireEntry()
		defer releaseEntry(entry)
		entry.Fields = l.fields

		entry.Data = buff.Bytes()
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

func (l *Logger) WithFields(fields Fields) *logger {
	return l.logger.WithFields(fields)
}

func (l *logger) Debug(msg ...interface{}) {
	l.output(DebugLevel, msg...)
}

func (l *logger) Info(msg ...interface{}) {
	l.output(InfoLevel, msg...)
}

func (l *logger) Warn(msg ...interface{}) {
	l.output(WarnLevel, msg...)
}

func (l *logger) Error(msg ...interface{}) {
	l.output(ErrorLevel, msg...)
}

func (l *logger) Fatal(msg ...interface{}) {
	l.output(FatalLevel, msg...)
}

func (l *logger) WithFields(fields Fields) *logger {
	log := acquireLogger()
	log.handlers = l.handlers
	for k, v := range fields {
		log.fields.Set(k, v)
	}
	return log
}

func (l *logger) Reset() {
	l.fields = l.fields[:0]
	l.handlers = l.handlers[:0]
}

package golog

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/millken/golog/internal/fasttime"
	"github.com/millken/x/buffer"
)

// var (
//
//	_ Logger = (*HandlerLog)(nil)
//
// )

type exitFunc func(int)

type Record struct {
	Config    Config
	Message   string
	Level     Level
	Time      time.Time
	calldepth int
	Fields    []Field
}

func (r Record) String() string {
	return fmt.Sprintf("%s %s %v", r.Level, r.Message, r.Fields)
}

type Logger struct {
	mu        sync.RWMutex // ensures atomic writes; protects the following fields
	fields    []Field
	calldepth int
	handlers  []Handler
	exit      exitFunc
}

func New(fields ...Field) *Logger {
	l := new(Logger)
	l.handlers = append(l.handlers, defaultHandler())
	l.exit = os.Exit

	l.setCalldepth(calldepth)
	// l.SetFields(fields...)
	// l.SetFlags(LstdFlags)

	return l
}

// AddHandler adds a handler to the logger
func (l *Logger) AddHandler(h Handler) {
	l.handlers = append(l.handlers, h)
}

// Fatalf calls underlying logger.Fatal.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.encodeOutput(FATAL, msg, nil)
	l.exit(1)
}

// Panicf calls underlying logger.Panic.
func (l *Logger) Panicf(format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.encodeOutput(PANIC, msg, nil)
	panic(msg)
}

// Debugf calls error log function if DEBUG level enabled.
func (l *Logger) Debugf(format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.encodeOutput(DEBUG, msg, nil)
}

// Infof calls error log function if INFO level enabled.
func (l *Logger) Infof(format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.encodeOutput(INFO, msg, nil)
}

// Warnf calls error log function if WARNING level enabled.
func (l *Logger) Warnf(format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.encodeOutput(WARNING, msg, nil)
}

// Errorf calls error log function if ERROR level enabled.
func (l *Logger) Errorf(format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.encodeOutput(ERROR, msg, nil)
}

// Fatal calls underlying logger.Fatal.
func (l *Logger) Fatal(msg string, keysAndVals ...interface{}) {
	l.encodeOutput(FATAL, msg, keysAndVals)
	l.exit(1)
}

// Panic calls underlying logger.Panic.
func (l *Logger) Panic(msg string, keysAndVals ...interface{}) {
	l.encodeOutput(PANIC, msg, keysAndVals)
	panic(l)
}

// Debug calls error log function if DEBUG level enabled.
func (l *Logger) Debug(msg string, keysAndVals ...interface{}) {
	l.encodeOutput(DEBUG, msg, keysAndVals)
}

// Info calls error log function if INFO level enabled.
func (l *Logger) Info(msg string, keysAndVals ...interface{}) {
	l.encodeOutput(INFO, msg, keysAndVals)
}

// Warn calls error log function if WARNING level enabled.
func (l *Logger) Warn(msg string, keysAndVals ...interface{}) {
	l.encodeOutput(WARNING, msg, keysAndVals)
}

// Error calls error log function if ERROR level enabled.
func (l *Logger) Error(msg string, keysAndVals ...interface{}) {
	l.encodeOutput(ERROR, msg, keysAndVals)
}

// WithValues returns a logger configured with the key-value pairs.
func (l *Logger) WithValues(keysAndVals ...interface{}) *Logger {
	clone := l.clone()
	for i := 0; i < len(keysAndVals); {
		if i == len(keysAndVals)-1 {
			break
		}
		key, val := keysAndVals[i], keysAndVals[i+1]
		keyStr, isString := key.(string)
		if !isString {
			break
		}

		clone.fields = append(clone.fields, Field{keyStr, val})
		i += 2
	}
	return clone
}

func (l *Logger) encodeOutput(level Level, msg string, args []interface{}) { //nolint:funlen
	rec := Record{
		Message:   msg,
		Level:     level,
		calldepth: l.calldepth,
		Time:      fasttime.Now(),
	}
	n := 0
	for _, f := range l.fields {
		rec.Fields = append(rec.Fields, f)
		n++
	}
	for i := 0; i < len(args); {
		if i == len(args)-1 {
			break
		}
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			break
		}

		rec.Fields = append(rec.Fields, Field{keyStr, val})
		n++
		i += 2
	}
	for _, h := range l.handlers {
		buf := buffer.Get()
		rec.Config = h.Config()
		if rec.Config.Level < level {
			continue
		}
		encoder := h.Encoder()
		if err := encoder.Encode(buf, rec); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to encode log, %v\n", err)
		}
		if _, err := h.Write(buf.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
		}
		buf.Free()
	}
}

func (l *Logger) setCalldepth(value int) {
	l.calldepth = value
}

// clone returns a copy of this "l" Logger.
func (l *Logger) clone() *Logger {
	return &Logger{
		fields:    l.fields,
		calldepth: l.calldepth,
		handlers:  l.handlers,
	}
}

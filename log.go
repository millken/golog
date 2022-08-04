package golog

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/millken/golog/internal/config"
	"github.com/millken/golog/internal/encoding"
	"github.com/millken/golog/internal/log"
	"github.com/millken/golog/internal/writer"
)

var (
	_ log.Logger = (*Log)(nil)
)

// Log is an implementation of Logger interface.
// It encapsulates default or custom logger to provide module and level based logging.
type Log struct {
	module        string
	fields        []log.Field
	once          sync.Once
	writer        io.Writer
	encoder       log.Encoder
	callerMap     map[log.Level]bool
	stacktraceMap map[log.Level]bool
	level         log.Level
}

// New creates and returns a Logger implementation based on given module name.
func NewLog(module string) *Log {
	callerMap := make(map[log.Level]bool, len(log.Levels))
	stacktraceMap := make(map[log.Level]bool, len(log.Levels))
	for _, v := range log.Levels {
		callerMap[v] = false
		stacktraceMap[v] = false
	}
	l := &Log{
		module:        module,
		callerMap:     callerMap,
		stacktraceMap: stacktraceMap,
	}
	l.init()
	return l
}

func (l *Log) init() {
	var err error
	l.once.Do(func() {
		mc := config.GetModuleConfig(l.module)
		switch mc.Writer.Type {
		case "file":
			l.writer, err = writer.NewFile(mc.Writer.FileConfig)
		default:
			l.writer, err = writer.NewFile(config.FileConfig{Path: "stdout"})
		}
		if err != nil {
			panic(err)
		}
		switch mc.Encoding {
		case "json":
		default:
			l.encoder = encoding.NewConsole(mc.ConsoleEncodingConfig)
		}
		l.level = mc.Level
		for _, v := range mc.CallerLevels {
			l.callerMap[v] = true
		}
		for _, v := range mc.StacktraceLevels {
			l.stacktraceMap[v] = true
		}
	})
}

// Fatalf calls underlying logger.Fatal.
func (l *Log) Fatalf(format string, args ...interface{}) {
	if l.level < log.FATAL {
		return
	}

	l.logf(log.FATAL, format, args...)
	os.Exit(1)
}

// Panicf calls underlying logger.Panic.
func (l *Log) Panicf(format string, args ...interface{}) {
	if l.level < log.PANIC {
		return
	}

	l.logf(log.PANIC, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Debugf calls error log function if DEBUG level enabled.
func (l *Log) Debugf(format string, args ...interface{}) {
	if l.level < log.DEBUG {
		return
	}

	l.logf(log.DEBUG, format, args...)
}

// Infof calls error log function if INFO level enabled.
func (l *Log) Infof(format string, args ...interface{}) {
	if l.level < log.INFO {
		return
	}

	l.logf(log.INFO, format, args...)
}

// Warnf calls error log function if WARNING level enabled.
func (l *Log) Warnf(format string, args ...interface{}) {
	if l.level < log.WARNING {
		return
	}

	l.logf(log.WARNING, format, args...)
}

// Errorf calls error log function if ERROR level enabled.
func (l *Log) Errorf(format string, args ...interface{}) {
	if l.level < log.ERROR {
		return
	}

	l.logf(log.ERROR, format, args...)
}

// Fatal calls underlying logger.Fatal.
func (l *Log) Fatal(msg string) {
	if l.level < log.FATAL {
		return
	}

	l.logf(log.FATAL, msg)
	os.Exit(1)
}

// Panic calls underlying logger.Panic.
func (l *Log) Panic(msg string) {
	if l.level < log.PANIC {
		return
	}

	l.logf(log.PANIC, msg)
	panic(msg)
}

// Debug calls error log function if DEBUG level enabled.
func (l *Log) Debug(msg string) {
	if l.level < log.DEBUG {
		return
	}

	l.logf(log.DEBUG, msg)
}

// Info calls error log function if INFO level enabled.
func (l *Log) Info(msg string) {
	if l.level < log.INFO {
		return
	}

	l.logf(log.INFO, msg)
}

// Warn calls error log function if WARNING level enabled.
func (l *Log) Warn(msg string) {
	if l.level < log.WARNING {
		return
	}

	l.logf(log.WARNING, msg)
}

// Error calls error log function if ERROR level enabled.
func (l *Log) Error(msg string) {
	if l.level < log.ERROR {
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
	e := log.AcquireEntry()
	e.Module = l.module

	copy(e.Fields[0:len(fields)], fields)
	e.SetFieldsLen(len(fields))

	e.Message = msg
	e.Level = level

	if config.IsCallerEnabled(e.Module, e.Level) {
		e.SetFlag(log.FlagCaller)
	}
	b, err := l.encoder.Encode(e)
	if err != nil {
		panic(err)
	}
	if _, err := l.writer.Write(b); err != nil {
		fmt.Fprintf(os.Stderr, "golog: failed to handle entry: %v", err)
	}
	e.Reset()
	log.ReleaseEntry(e)
}

// LoadConfig - Load configuration from file
func LoadConfig(path string) error {
	return config.Load(path)
}

//nolint:gochecknoglobals
var (
	loggerProviderInstance log.Logger
	loggerProviderOnce     sync.Once
)

// Panic logs a message using Panic level and panics.
func Panic(msg string) {
	loggerProvider().Panic(msg)
}

// Fatal logs a message using Fatal level and exits with status 1.
func Fatal(msg string) {
	loggerProvider().Fatal(msg)
}

// Error logs a message using Error level.
func Error(msg string) {
	loggerProvider().Error(msg)
}

// Warn logs a message using Warn level.
func Warn(msg string) {
	loggerProvider().Warn(msg)
}

// Info logs a message using Info level.
func Info(msg string) {
	loggerProvider().Infof(msg)
}

// Debug logs a message using Debug level.
func Debug(msg string) {
	loggerProvider().Debugf(msg)
}

func loggerProvider() log.Logger {
	loggerProviderOnce.Do(func() {
		loggerProviderInstance = NewLog("")
	})

	return loggerProviderInstance
}

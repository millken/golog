package golog

import (
	"fmt"
	"io"
	"os"
	"sync"
)

var (
	_ Logger = (*Log)(nil)
)

// Log is an implementation of Logger interface.
// It encapsulates default or custom logger to provide module and level based logging.
type Log struct {
	module        string
	fields        []Field
	once          sync.Once
	writer        io.Writer
	encoder       Encoder
	callerMap     map[Level]bool
	callerSkip    int
	stacktraceMap map[Level]bool
	level         Level
}

func newLogger() *Log {
	callerMap := make(map[Level]bool, len(Levels))
	stacktraceMap := make(map[Level]bool, len(Levels))
	for _, v := range Levels {
		callerMap[v] = false
		stacktraceMap[v] = false
	}
	return &Log{
		callerMap:     callerMap,
		stacktraceMap: stacktraceMap,
	}
}

// New creates and returns a Logger implementation based on given module name.
func New(module string) *Log {
	l := newLogger()
	l.module = module
	l.init()
	return l
}

// NewLoggerByConfig creates and returns a Logger implementation based on given config.
func NewLoggerByConfig(module string, cfg Config) (*Log, error) {
	l := newLogger()
	l.module = module
	err := l.initConfig(cfg)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (l *Log) init() {
	var err error
	l.once.Do(func() {
		mc := GetModuleConfig(l.module)
		if err = l.initConfig(mc); err != nil {
			panic(err)
		}
	})
}

func (l *Log) initConfig(cfg Config) error {
	var err error
	switch cfg.Writer.Type {
	case "file":
		l.writer, err = NewFile(cfg.Writer.FileConfig)
	case "custom":
		l.writer = cfg.Writer.CustomWriter
	default:
		l.writer, err = NewFile(FileConfig{Path: "stdout"})
	}
	if err != nil {
		return err
	}
	switch cfg.Encoding {
	case "json":
		l.encoder = NewJSONEncoder(cfg.JSONEncoderConfig)
	default:
		l.encoder = NewConsoleEncoder(cfg.ConsoleEncoderConfig)
	}
	l.level = DefaultLevel // if level is not set, set it to INFO
	if cfg.Level > 0 {
		l.level = cfg.Level
	}
	for _, v := range cfg.CallerLevels {
		l.callerMap[v] = true
	}
	for _, v := range cfg.StacktraceLevels {
		l.stacktraceMap[v] = true
	}
	return nil
}

//CallerSkip is used to set the number of caller frames to skip.
func (l *Log) CallerSkip(skip int) *Log {
	l.callerSkip += skip
	return l
}

// Fatalf calls underlying logger.Fatal.
func (l *Log) Fatalf(format string, args ...interface{}) {
	if l.level < FATAL {
		return
	}

	l.logf(FATAL, format, args...)
	os.Exit(1)
}

// Panicf calls underlying logger.Panic.
func (l *Log) Panicf(format string, args ...interface{}) {
	if l.level < PANIC {
		return
	}

	l.logf(PANIC, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Debugf calls error log function if DEBUG level enabled.
func (l *Log) Debugf(format string, args ...interface{}) {
	if l.level < DEBUG {
		return
	}

	l.logf(DEBUG, format, args...)
}

// Infof calls error log function if INFO level enabled.
func (l *Log) Infof(format string, args ...interface{}) {
	if l.level < INFO {
		return
	}

	l.logf(INFO, format, args...)
}

// Warnf calls error log function if WARNING level enabled.
func (l *Log) Warnf(format string, args ...interface{}) {
	if l.level < WARNING {
		return
	}

	l.logf(WARNING, format, args...)
}

// Errorf calls error log function if ERROR level enabled.
func (l *Log) Errorf(format string, args ...interface{}) {
	if l.level < ERROR {
		return
	}

	l.logf(ERROR, format, args...)
}

// Fatal calls underlying logger.Fatal.
func (l *Log) Fatal(msg string, fields ...Field) {
	if l.level < FATAL {
		return
	}

	l.log(FATAL, msg, fields)
	os.Exit(1)
}

// Panic calls underlying logger.Panic.
func (l *Log) Panic(msg string, fields ...Field) {
	if l.level < PANIC {
		return
	}

	l.log(PANIC, msg, fields)
	panic(msg)
}

// Debug calls error log function if DEBUG level enabled.
func (l *Log) Debug(msg string, fields ...Field) {
	if l.level < DEBUG {
		return
	}

	l.log(DEBUG, msg, fields)
}

// Info calls error log function if INFO level enabled.
func (l *Log) Info(msg string, field ...Field) {
	if l.level < INFO {
		return
	}

	l.log(INFO, msg, field)
}

// Warn calls error log function if WARNING level enabled.
func (l *Log) Warn(msg string, field ...Field) {
	if l.level < WARNING {
		return
	}

	l.log(WARNING, msg, field)
}

// Error calls error log function if ERROR level enabled.
func (l *Log) Error(msg string, field ...Field) {
	if l.level < ERROR {
		return
	}

	l.log(ERROR, msg, field)
}

// WithField returns a logger configured with the key-value pair.
func (l *Log) WithField(k string, v interface{}) Logger {
	clone := l.Clone()
	clone.fields = append(l.fields, Field{Key: k, Val: v})
	return clone
}

// WithFields returns a logger configured with the key-value pairs.
func (l *Log) WithFields(fields Fields) Logger {
	clone := l.Clone()
	for k, v := range fields {
		clone.fields = append(l.fields, Field{Key: k, Val: v})
	}
	return clone
}

func (l *Log) logf(level Level, format string, args ...interface{}) {
	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.output(level, msg, l.fields)
}

func (l *Log) log(level Level, msg string, fields []Field) {
	if len(l.fields) > 0 {
		fields = append(l.fields, fields...)
	}
	l.output(level, msg, fields)
}

func (l *Log) output(level Level, msg string, fields []Field) {
	e := acquireEntry()
	defer releaseEntry(e)
	e.Module = l.module

	e.Fields = fields
	e.SetFieldsLen(len(fields))

	e.Message = msg
	e.Level = level
	e.SetCallerSkip(l.callerSkip)

	if l.isCallerEnabled(e.Level) {
		e.SetFlag(FlagCaller)
	}
	if l.isStacktraceEnabled(e.Level) {
		e.SetFlag(FlagStacktrace)
	}
	b, err := l.encoder.Encode(e)
	if err != nil {
		panic(err)
	}
	if _, err := l.writer.Write(b); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write log: %v", err)
	}
	e.Reset()
}

func (l *Log) isCallerEnabled(level Level) bool {
	if enabled, ok := l.callerMap[level]; ok {
		return enabled
	}
	return false
}

func (l *Log) isStacktraceEnabled(level Level) bool {
	if enabled, ok := l.stacktraceMap[level]; ok {
		return enabled
	}
	return false
}

// Clone returns a copy of this "l" Logger.
// This copy is returned as pointer as well.
func (l *Log) Clone() *Log {
	copy := *l
	copy.callerSkip = 0
	return &copy
}

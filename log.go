package golog

import (
	"fmt"
	"io"
	"os"
	"slices"
	"sync"
)

var (
	_ Logger = (*Log)(nil)
)

// Log is an implementation of Logger interface.
// It encapsulates default or custom logger to provide module and level based logging.
type Log struct {
	module     string
	fields     []Field
	once       sync.Once
	writer     io.Writer
	encoder    Encoder
	callerLvl  uint32
	callerSkip int
	tracerLvl  uint32
	level      Level
}

func newLogger() *Log {
	return &Log{
		fields: []Field{},
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
	l.once.Do(func() {
		mc := GetModuleConfig(l.module)
		if err := l.initConfig(mc); err != nil {
			panic(err)
		}
	})
}

func (l *Log) initConfig(cfg Config) error {
	var err error
	switch cfg.Handler.Type {
	case HandlerTypeFile:
		l.writer, err = NewFile(cfg.Handler.File)
	case HandlerTypeRotateFile:
		l.writer, err = NewRotateFile(cfg.Handler.RotateFile)
	case HandlerTypeCustom:
		l.writer = cfg.Handler.Writer
	default:
		l.writer, err = NewFile(FileConfig{Path: "stdout"})
	}
	if err != nil {
		return err
	}
	switch cfg.Encoding {
	case JSONEncoding:
		l.encoder = NewJSONEncoder(cfg.JSONEncoder)
	default:
		l.encoder = NewTextEncoder(cfg.TextEncoder)
	}
	l.level = INFO // if level is not set, set it to INFO
	if cfg.Level > 0 {
		l.level = cfg.Level
	}
	for _, v := range cfg.CallerLevels {
		l.callerLvl |= uint32(v)
	}
	for _, v := range cfg.StacktraceLevels {
		l.tracerLvl |= uint32(v)
	}
	return nil
}

// CallerSkip is used to set the number of caller frames to skip.
func (l *Log) CallerSkip(skip int) *Log {
	l.callerSkip = skip
	return l
}

func formatMessage(format string, args ...any) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

func (l *Log) logf(level Level, format string, args ...any) {
	if l.level < level {
		return
	}
	msg := formatMessage(format, args...)
	l.output(level, msg, nil, 1)
}

// Fatalf calls underlying logger.Fatal.
func (l *Log) Fatalf(format string, args ...any) {
	if l.level < FATAL {
		return
	}
	msg := formatMessage(format, args...)
	l.output(FATAL, msg, nil, 1)
	os.Exit(1)
}

// Panicf calls underlying logger.Panic.
func (l *Log) Panicf(format string, args ...any) {
	if l.level < PANIC {
		return
	}
	msg := formatMessage(format, args...)
	l.output(PANIC, msg, nil, 1)
	panic(msg)
}

// Debugf calls debug log function if DEBUG level enabled.
func (l *Log) Debugf(format string, args ...any) {
	l.logf(DEBUG, format, args...)
}

// Infof calls info log function if INFO level enabled.
func (l *Log) Infof(format string, args ...any) {
	l.logf(INFO, format, args...)
}

// Warnf calls warn log function if WARNING level enabled.
func (l *Log) Warnf(format string, args ...any) {
	l.logf(WARNING, format, args...)
}

// Errorf calls error log function if ERROR level enabled.
func (l *Log) Errorf(format string, args ...any) {
	l.logf(ERROR, format, args...)
}

// Fatal calls underlying logger.Fatal.
func (l *Log) Fatal(msg string, keysAndVals ...any) {
	if l.level < FATAL {
		return
	}

	l.output(FATAL, msg, keysAndVals, 0)
	os.Exit(1)
}

// Panic calls underlying logger.Panic.
func (l *Log) Panic(msg string, keysAndVals ...any) {
	if l.level < PANIC {
		return
	}

	l.output(PANIC, msg, keysAndVals, 0)
	panic(msg)
}

// Debug calls debug log function if DEBUG level enabled.
func (l *Log) Debug(msg string, keysAndVals ...any) {
	if l.level < DEBUG {
		return
	}

	l.output(DEBUG, msg, keysAndVals, 0)
}

// Info calls info log function if INFO level enabled.
func (l *Log) Info(msg string, keysAndVals ...any) {
	if l.level < INFO {
		return
	}

	l.output(INFO, msg, keysAndVals, 0)
}

// Warn calls warn log function if WARNING level enabled.
func (l *Log) Warn(msg string, keysAndVals ...any) {
	if l.level < WARNING {
		return
	}

	l.output(WARNING, msg, keysAndVals, 0)
}

// Error calls error log function if ERROR level enabled.
func (l *Log) Error(msg string, keysAndVals ...any) {
	if l.level < ERROR {
		return
	}
	l.output(ERROR, msg, keysAndVals, 0)
}

// WithValues returns a logger configured with the key-value pairs.
func (l *Log) WithValues(keysAndVals ...any) Logger {
	clone := l.clone()
	for i := 0; i+1 < len(keysAndVals); i += 2 {
		key, val := keysAndVals[i], keysAndVals[i+1]
		keyStr, isString := key.(string)
		if !isString {
			fmt.Fprintf(os.Stderr, "golog: WithValues received non-string key: %v, ignoring remaining args\n", key)
			break
		}

		clone.fields = append(clone.fields, field(keyStr, val))
	}
	if len(keysAndVals)%2 != 0 {
		fmt.Fprintf(os.Stderr, "golog: WithValues received odd number of arguments, ignoring last key: %v\n", keysAndVals[len(keysAndVals)-1])
	}
	return clone
}

func (l *Log) output(level Level, msg string, args []any, extraCallerSkip int) { //nolint:funlen
	e := acquireEntry()
	defer releaseEntry(e)
	e.Module = l.module

	n := 0
	for _, f := range l.fields {
		e.Fields = append(e.Fields, f)
		n++
	}
	for i := 0; i+1 < len(args); i += 2 {
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			break
		}

		e.Fields = append(e.Fields, field(keyStr, val))
		n++
	}
	e.SetFieldsLen(n)

	e.Message = msg
	e.Level = level
	e.SetCallerSkip(l.callerSkip + extraCallerSkip)

	if l.isCallerEnabled(e.Level) {
		e.SetFlag(FlagCaller)
	}
	if l.isStacktraceEnabled(e.Level) {
		e.SetFlag(FlagStacktrace)
	}
	b, err := l.encoder.Encode(e)
	if err != nil {
		fmt.Fprintf(os.Stderr, "golog: failed to encode log: %v\n", err)
		return
	}
	if _, err := l.writer.Write(b); err != nil {
		fmt.Fprintf(os.Stderr, "golog: failed to write log: %v\n", err)
	}
}

func (l *Log) isCallerEnabled(level Level) bool {
	return l.callerLvl&uint32(level) == uint32(level)
}

func (l *Log) isStacktraceEnabled(level Level) bool {
	return l.tracerLvl&uint32(level) == uint32(level)
}

// clone returns a copy of this "l" Logger.
func (l *Log) clone() *Log {
	fields := slices.Clone(l.fields)
	return &Log{
		level:      l.level,
		module:     l.module,
		writer:     l.writer,
		fields:     fields,
		encoder:    l.encoder,
		callerLvl:  l.callerLvl,
		callerSkip: l.callerSkip,
		tracerLvl:  l.tracerLvl,
		once:       sync.Once{},
	}
}

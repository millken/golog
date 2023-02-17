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
		fields:        []Field{},
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
	l.level = INFO // if level is not set, set it to INFO
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

	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.output(FATAL, msg, nil)
	os.Exit(1)
}

// Panicf calls underlying logger.Panic.
func (l *Log) Panicf(format string, args ...interface{}) {
	if l.level < PANIC {
		return
	}

	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.output(PANIC, msg, nil)
	panic(msg)
}

// Debugf calls error log function if DEBUG level enabled.
func (l *Log) Debugf(format string, args ...interface{}) {
	if l.level < DEBUG {
		return
	}

	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.output(DEBUG, msg, nil)
}

// Infof calls error log function if INFO level enabled.
func (l *Log) Infof(format string, args ...interface{}) {
	if l.level < INFO {
		return
	}

	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.output(INFO, msg, nil)
}

// Warnf calls error log function if WARNING level enabled.
func (l *Log) Warnf(format string, args ...interface{}) {
	if l.level < WARNING {
		return
	}

	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.output(WARNING, msg, nil)
}

// Errorf calls error log function if ERROR level enabled.
func (l *Log) Errorf(format string, args ...interface{}) {
	if l.level < ERROR {
		return
	}

	var msg string
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	} else {
		msg = format
	}
	l.output(ERROR, msg, nil)
}

// Fatal calls underlying logger.Fatal.
func (l *Log) Fatal(msg string, keysAndVals ...interface{}) {
	if l.level < FATAL {
		return
	}

	l.output(FATAL, msg, keysAndVals)
	os.Exit(1)
}

// Panic calls underlying logger.Panic.
func (l *Log) Panic(msg string, keysAndVals ...interface{}) {
	if l.level < PANIC {
		return
	}

	l.output(PANIC, msg, keysAndVals)
	panic(msg)
}

// Debug calls error log function if DEBUG level enabled.
func (l *Log) Debug(msg string, keysAndVals ...interface{}) {
	if l.level < DEBUG {
		return
	}

	l.output(DEBUG, msg, keysAndVals)
}

// Info calls error log function if INFO level enabled.
func (l *Log) Info(msg string, keysAndVals ...interface{}) {
	if l.level < INFO {
		return
	}

	l.output(INFO, msg, keysAndVals)
}

// Warn calls error log function if WARNING level enabled.
func (l *Log) Warn(msg string, keysAndVals ...interface{}) {
	if l.level < WARNING {
		return
	}

	l.output(WARNING, msg, keysAndVals)
}

// Error calls error log function if ERROR level enabled.
func (l *Log) Error(msg string, keysAndVals ...interface{}) {
	if l.level < ERROR {
		return
	}
	l.output(ERROR, msg, keysAndVals)
}

// WithValues returns a logger configured with the key-value pairs.
func (l *Log) WithValues(keysAndVals ...interface{}) Logger {
	clone := l.clone()
	for i := 0; i < len(keysAndVals); {
		if fields, ok := keysAndVals[i].(Fields); ok {
			for k, v := range fields {
				clone.fields = append(clone.fields, field(k, v))
			}
			i++
			continue
		} else if field, ok := keysAndVals[i].(Field); ok {
			clone.fields = append(clone.fields, field)
			i++
			continue
		}
		if i == len(keysAndVals)-1 {
			break
		}
		key, val := keysAndVals[i], keysAndVals[i+1]
		keyStr, isString := key.(string)
		if !isString {
			break
		}

		clone.fields = append(clone.fields, field(keyStr, val))
		i += 2
	}
	return clone
}

func (l *Log) output(level Level, msg string, args []interface{}) {
	e := acquireEntry()
	defer releaseEntry(e)
	e.Module = l.module

	n := 0
	for _, f := range l.fields {
		e.Fields = append(e.Fields, f)
		n++
	}
	for i := 0; i < len(args); {
		if fields, ok := args[i].(Fields); ok {
			for k, v := range fields {
				e.Fields = append(e.Fields, field(k, v))
				n++
			}
			i++
			continue
		} else if field, ok := args[i].(Field); ok {
			e.Fields = append(e.Fields, field)
			i++
			n++
			continue
		}
		if i == len(args)-1 {
			break
		}
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			break
		}

		e.Fields = append(e.Fields, field(keyStr, val))
		n++
		i += 2
	}
	//e.Fields = append(e.Fields, l.handleFields(args)...)
	e.SetFieldsLen(n)

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

// clone returns a copy of this "l" Logger.
func (l *Log) clone() *Log {
	return &Log{
		level:         l.level,
		module:        l.module,
		writer:        l.writer,
		fields:        l.fields,
		encoder:       l.encoder,
		callerMap:     l.callerMap,
		stacktraceMap: l.stacktraceMap,
		once:          sync.Once{},
	}
}

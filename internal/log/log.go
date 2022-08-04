package log

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Level is a log level for a logging message.
type Level int

// Log levels.
const (
	PANIC Level = iota
	FATAL
	ERROR
	WARNING
	INFO
	DEBUG
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case PANIC:
		return "panic"
	case FATAL:
		return "fatal"
	case ERROR:
		return "error"
	case WARNING:
		return "warning"
	case INFO:
		return "info"
	case DEBUG:
		return "debug"
	default:
		return "unknown"
	}
}

func (l Level) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(l.String())
}

func (l *Level) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var v string
	if err := unmarshal(&v); err != nil {
		return err
	}
	level, err := ParseLevel(v)
	if err != nil {
		return err
	}
	*l = level
	return nil
}

func ParseLevel(level string) (Level, error) {
	var l Level
	switch strings.ToLower(level) {
	case "panic":
		l = PANIC
	case "fatal":
		l = FATAL
	case "error":
		l = ERROR
	case "warning", "warn":
		l = WARNING
	case "info":
		l = INFO
	case "debug":
		l = DEBUG
	default:
		return l, fmt.Errorf("unknown log level: %s", level)
	}
	return l, nil
}

const (
	// TimestampFieldName is the field name used for the timestamp field.
	TimestampFieldName = "time"

	// LevelFieldName is the field name used for the level field.
	LevelFieldName = "level"

	// MessageFieldName is the field name used for the message field.
	MessageFieldName = "message"

	// ErrorFieldName is the field name used for error fields.
	ErrorFieldName = "error"

	// CallerFieldName is the field name used for caller field.
	CallerFieldName = "caller"

	// ErrorStackFieldName is the field name used for error stacks.
	ErrorStackFieldName = "stack"
)

// Field is a key/value pair.
type Field struct {
	Key string
	Val interface{}
}

// Logger represents a general-purpose logger.
type Logger interface {
	WithField(k string, v interface{}) Logger
	WithFields(fields ...Field) Logger
	Panicf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Debugf(msg string, args ...interface{})
	Panic(msg string)
	Fatal(msg string)
	Error(msg string)
	Warn(msg string)
	Info(msg string)
	Debug(msg string)
}

// LoggerProvider is a factory for moduled loggers.
type LoggerProvider interface {
	GetLogger(module string) Logger
}

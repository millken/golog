package golog

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	// Levels is a slice of all log levels.
	Levels = []Level{
		PANIC,
		FATAL,
		ERROR,
		WARNING,
		INFO,
		DEBUG,
	}
)

// Level is a log level for a logging message.
type Level uint32

// Log levels.
const (
	PANIC Level = 1 << iota
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
	}
	return "unknown"
}

// MarshalYAML implements the yaml.Marshaler interface.
func (l Level) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(l.String())
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
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

// ParseLevel parses a string into a log level.
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

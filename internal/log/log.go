package log

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

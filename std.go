package golog

import (
	"io"
	"os"
)

var std = newStd()

type stdlog struct {
	handler Handler
	logger  *logger
}

func newStd() *stdlog {
	handler := &FileHandler{
		Output: os.Stderr,
	}
	handler.SetLevel(InfoLevel)
	handler.SetFormatter(&TextFormatter{
		CallerSkipFrameCount: 7,
		EnableCaller:         true,
	})
	logger := newLogger()
	logger.AddHandler(handler)
	return &stdlog{
		handler: handler,
		logger:  logger,
	}
}

func SetLevel(level Level) {
	std.handler.SetLevel(level)
}

func SetOutput(output io.Writer) {
	std.handler.(*FileHandler).SetOutput(output)
}

func Fatal(msg string, fields ...field) {
	std.logger.Fatal(msg, fields...)
}

func Error(msg string, fields ...field) {
	std.logger.Error(msg, fields...)
}

func Warn(msg string, fields ...field) {
	std.logger.Warn(msg, fields...)
}

func Info(msg string, fields ...field) {
	std.logger.Info(msg, fields...)
}

func Debug(msg string, fields ...field) {
	std.logger.Debug(msg, fields...)
}

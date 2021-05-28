package golog

import (
	"io"
	"os"
	"sync"
)

var (
	stdOnce         sync.Once
	stdHandler      Handler
	stdLogger       *Logger
	stdFormatter    Formatter
	StdTimeFormat   = "2006-01-02 15:04:05"
	StdNoColor      = true
	StdEnableCaller = false
)

func init() {
	stdOnce.Do(func() {
		stdHandler = &FileHandler{
			Output: os.Stderr,
		}
		stdHandler.SetLevel(InfoLevel)
		stdFormatter = &TextFormatter{
			NoColor:              StdNoColor,
			TimeFormat:           StdTimeFormat,
			CallerSkipFrameCount: 6,
			EnableCaller:         StdEnableCaller,
		}
		stdHandler.SetFormatter(stdFormatter)
		stdLogger = NewLogger()
		stdLogger.AddHandler(stdHandler)
	})
}

func SetLevel(level Level) {
	stdHandler.SetLevel(level)
}

func SetOutput(output io.Writer) {
	stdHandler.(*FileHandler).SetOutput(output)
}

func Fatal(msg string, fields ...field) {
	stdLogger.Fatal(msg, fields...)
}

func Error(msg string, fields ...field) {
	stdLogger.Error(msg, fields...)
}

func Warn(msg string, fields ...field) {
	stdLogger.Warn(msg, fields...)
}

func Info(msg string, fields ...field) {
	stdLogger.Info(msg, fields...)
}

func Debug(msg string, fields ...field) {
	stdLogger.Debug(msg, fields...)
}

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
	stdTimeFormat   = "2006-01-02 15:04:05"
	stdNoColor      = false
	stdEnableCaller = false
)

func init() {
	stdOnce.Do(func() {
		stdHandler = &FileHandler{
			Output: os.Stderr,
		}
		stdHandler.SetLevel(InfoLevel)
		stdFormatter = &TextFormatter{
			NoColor:              stdNoColor,
			TimeFormat:           stdTimeFormat,
			CallerSkipFrameCount: 6,
			EnableCaller:         stdEnableCaller,
		}
		stdHandler.SetFormatter(stdFormatter)
		stdLogger = NewLogger()
		stdLogger.AddHandler(stdHandler)
	})
}

func StdSetLevel(level Level) {
	stdHandler.SetLevel(level)
}

func StdSetOutput(output io.Writer) {
	stdHandler.(*FileHandler).SetOutput(output)
}

func StdEnableCaller() {
	stdHandler.(*FileHandler).GetFormatter().(*TextFormatter).EnableCaller = true
}
func StdNoColor() {
	stdHandler.(*FileHandler).GetFormatter().(*TextFormatter).NoColor = true
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

package golog

import (
	"io"
	"os"
)

var (
	stdTimeFormat = "2006-01-02 15:04:05"
)

type stdLogger struct {
	handler Handler
	*Logger
}

func NewStdLog() *stdLogger {
	stdHandler := &FileHandler{
		Output: os.Stderr,
	}
	stdHandler.SetLevel(InfoLevel)
	stdFormatter := &TextFormatter{
		NoColor:              false,
		TimeFormat:           stdTimeFormat,
		CallerSkipFrameCount: 6,
		EnableCaller:         false,
	}
	stdHandler.SetFormatter(stdFormatter)

	log := NewLogger()
	log.AddHandler(stdHandler)
	return &stdLogger{
		handler: stdHandler,
		Logger:  log,
	}
}

func (l *stdLogger) SetLevel(level Level) {
	l.handler.SetLevel(level)
}

func (l *stdLogger) SetOutput(output io.Writer) {
	l.handler.(*FileHandler).SetOutput(output)
}

func (l *stdLogger) EnableCaller(enable bool) {
	l.handler.(*FileHandler).Formatter().(*TextFormatter).EnableCaller = enable
}

func (l *stdLogger) EnableColor(enable bool) {
	l.handler.(*FileHandler).Formatter().(*TextFormatter).NoColor = !enable
}

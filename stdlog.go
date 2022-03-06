package golog

import (
	"io"
	"os"
)

var (
	StdTimeFormat = "2006-01-02 15:04:05"
)

type StdOption struct {
	Level                Level
	Output               io.Writer
	NoColor              bool
	TimeFormat           string
	CallerSkipFrameCount int
	EnableCaller         bool
	DisableTimestamp     bool
	PartsOrder           []string
	PartsExclude         []string
}

func prepareStdOptions(opt StdOption) StdOption {
	if opt.Output == nil {
		opt.Output = os.Stderr
	}
	if opt.TimeFormat == "" {
		opt.TimeFormat = consoleDefaultTimeFormat
	}
	if opt.Level == NoLevel {
		opt.Level = DefaultLevel
	}
	if opt.CallerSkipFrameCount == 0 {
		opt.CallerSkipFrameCount = 6
	}
	if opt.PartsOrder == nil {
		opt.PartsOrder = consoleDefaultPartsOrder()
	}
	return opt
}

func NewStdLog(opts ...StdOption) *Logger {
	var o StdOption
	if len(opts) > 0 {
		o = opts[0]
	}
	opt := prepareStdOptions(o)
	stdHandler := &FileHandler{
		Output: opt.Output,
	}
	stdHandler.SetLevel(opt.Level)

	stdFormatter := &TextFormatter{
		NoColor:              opt.NoColor,
		EnableCaller:         opt.EnableCaller,
		CallerSkipFrameCount: opt.CallerSkipFrameCount,
		TimeFormat:           opt.TimeFormat,
		DisableTimestamp:     opt.DisableTimestamp,
		PartsOrder:           opt.PartsOrder,
		PartsExclude:         opt.PartsExclude,
	}
	stdHandler.SetFormatter(stdFormatter)

	logger := NewLogger()
	logger.AddHandler(stdHandler)
	return logger
}

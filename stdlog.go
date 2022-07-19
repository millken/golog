package golog

import (
	"io"
	"os"
)

var (
	// StdTimeFormat is the default time format for the console logger.
	StdTimeFormat = "2006-01-02 15:04:05"
)

// StdOption is the option for StdLog.
type StdOption struct {
	NoColor          bool
	EnableCaller     bool
	DisableTimestamp bool
	Level            Level
	Output           io.Writer
	TimeFormat       string
	PartsOrder       []string
	PartsExclude     []string
}

func prepareStdOptions(opt *StdOption) *StdOption {
	if opt.Output == nil {
		opt.Output = os.Stderr
	}
	if opt.TimeFormat == "" {
		opt.TimeFormat = consoleDefaultTimeFormat
	}
	if opt.PartsOrder == nil {
		opt.PartsOrder = consoleDefaultPartsOrder()
	}
	return opt
}

// NewStdLog creates a new StdLog.
func NewStdLog(opts ...StdOption) *Logger {
	var o StdOption
	if len(opts) > 0 {
		o = opts[0]
	}
	opt := prepareStdOptions(&o)
	stdHandler := NewLoggerHandler(opt.Output)

	stdHandler.SetLevel(opt.Level)

	stdFormatter := &TextFormatter{
		NoColor:          opt.NoColor,
		EnableCaller:     opt.EnableCaller,
		TimeFormat:       opt.TimeFormat,
		DisableTimestamp: opt.DisableTimestamp,
		PartsOrder:       opt.PartsOrder,
		PartsExclude:     opt.PartsExclude,
	}
	stdHandler.SetFormatter(stdFormatter)

	logger := NewLogger()
	logger.AddHandler(stdHandler)
	return logger
}

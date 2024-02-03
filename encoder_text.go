package golog

import (
	"github.com/millken/golog/internal/stack"
	"github.com/millken/x/buffer"
)

// TimestampFormat typenc.
type TimestampFormat int

// Logger timestamp formats.
const (
	TimestampFormatSeconds TimestampFormat = iota + 1
	TimestampFormatNanoseconds
)

const (
	calldepth    = 3
	calldepthStd = calldepth + 1
)

const unknownFile = "???"

const (
	printLevelStr   = ""
	panicLevelStr   = "PANIC"
	fatalLevelStr   = "FATAL"
	errorLevelStr   = "ERROR"
	warningLevelStr = "WARNING"
	infoLevelStr    = "INFO"
	debugLevelStr   = "DEBUG"
	traceLevelStr   = "TRACE"
)

// EncoderText is the text enconder.
type EncoderText struct {
	cfg EncoderTextConfig
}

// EncoderTextConfig is the configuration of text encoder.
type EncoderTextConfig struct {
	// Default: -
	Separator string

	// Default: timenc.RFC3339
	DatetimeLayout string

	// Default: TimestampFormatSeconds
	TimestampFormat TimestampFormat

	// TimeFormat specifies the format for timestamp in output.
	TimeFormat string `json:"timeFormat" yaml:"timeFormat"`
	// DisableTimestamp disables the timestamp in output.
	DisableTimestamp bool `json:"disableTimestamp" yaml:"disableTimestamp"`
	// CallerSkipFrame is the number of stack frames to skip when reporting the calling function.
	CallerSkipFrame int `json:"callerSkipFrame" yaml:"callerSkipFrame"`
}

// NewEncoderText creates a new text encoder.
func NewEncoderText(cfg EncoderTextConfig) *EncoderText {
	if cfg.Separator == "" {
		cfg.Separator = defaultTextSeparator
	}

	if cfg.DatetimeLayout == "" {
		cfg.DatetimeLayout = defaultDatetimeLayout
	}

	if cfg.TimestampFormat == 0 {
		cfg.TimestampFormat = defaultTimestampFormat
	}

	enc := new(EncoderText)
	enc.cfg = cfg

	return enc
}

// Encode encodes the given entry to the buffer.
func (enc *EncoderText) Encode(buf *buffer.Buffer, e Record) error {
	if e.Config.UTC {
		e.Time = e.Time.UTC()
	}
	if e.Config.Datetime {
		buf.AppendTime(e.Time, enc.cfg.DatetimeLayout)
		buf.WriteString(enc.cfg.Separator) // nolint:errcheck
	}

	if e.Config.Timestamp {
		switch enc.cfg.TimestampFormat {
		case TimestampFormatSeconds:
			buf.AppendInt(e.Time.Unix())
		case TimestampFormatNanoseconds:
			buf.AppendInt(e.Time.UnixNano())
		}
		buf.WriteString(enc.cfg.Separator) // nolint:errcheck
	}

	if levelStr := e.Level.String(); levelStr != "" {
		buf.WriteString(levelStr)          // nolint:errcheck
		buf.WriteString(enc.cfg.Separator) // nolint:errcheck
	}

	if e.Config.Shortfile || e.Config.Longfile || e.Config.Stack {
		stackSkip := enc.cfg.CallerSkipFrame + e.calldepth
		frames := stack.Tracer(stackSkip, e.Config.Stack)

		if len(frames) > 0 {
			if e.Config.Shortfile || e.Config.Longfile {
				file := frames[0].File

				if e.Config.Shortfile {
					for i := len(file) - 1; i > 0; i-- {
						if file[i] == '/' {
							file = file[i+1:]
							break
						}
					}
				}
				buf.WriteString(file)                // nolint:errcheck
				buf.WriteByte(':')                   // nolint:errcheck
				buf.AppendInt(int64(frames[0].Line)) // nolint:errcheck
				buf.WriteString(enc.cfg.Separator)
			}
			if e.Config.Function {
				buf.WriteString(frames[0].Function) // nolint:errcheck
				buf.WriteString(enc.cfg.Separator)  // nolint:errcheck
			}
			if e.Config.Stack {
				buf.WriteNewLine()
				stackfmt := stack.NewStackFormatter(buf)
				stackfmt.FormatFrames(frames)
				buf.WriteNewLine()
			}
		}
	}

	for _, field := range e.Fields {
		buf.WriteString(field.Key) // nolint:errcheck
		buf.WriteByte('=')         // nolint:errcheck
		buf.WriteInterface(field.Val)
		buf.WriteString(enc.cfg.Separator) // nolint:errcheck
	}

	buf.WriteString(e.Message) // nolint:errcheck
	buf.WriteNewLine()

	return nil
}

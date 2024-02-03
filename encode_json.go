package golog

import (
	"strconv"

	"github.com/millken/golog/internal/stack"
	"github.com/millken/x/buffer"
)

// EncoderJSON is the json encoder.
type EncoderJSON struct {
	cfg EncoderJSONConfig
}

// EncoderJSONConfig is the configuration of json encoder.
type EncoderJSONConfig struct {
	FieldMap EnconderJSONFieldMap

	// Default: time.RFC3339
	DatetimeLayout string

	// Default: TimestampFormatSeconds
	TimestampFormat TimestampFormat
}

// EnconderJSONFieldMap defines name of keys.
type EnconderJSONFieldMap struct {
	// Default: datetime
	DatetimeKey string

	// Default: timestamp
	TimestampKey string

	// Default: level
	LevelKey string

	// Default: file
	FileKey string

	// Default: stack
	StackKey string

	// Default: message
	MessageKey string
}

// NewEncoderJSON creates a new json encoder.
func NewEncoderJSON(cfg EncoderJSONConfig) *EncoderJSON {
	if cfg.FieldMap.DatetimeKey == "" {
		cfg.FieldMap.DatetimeKey = defaultJSONFieldKeyDatetime
	}

	if cfg.FieldMap.TimestampKey == "" {
		cfg.FieldMap.TimestampKey = defaultJSONFieldKeyTimestamp
	}

	if cfg.FieldMap.LevelKey == "" {
		cfg.FieldMap.LevelKey = defaultJSONFieldKeyLevel
	}

	if cfg.FieldMap.FileKey == "" {
		cfg.FieldMap.FileKey = defaultJSONFieldKeyFile
	}

	if cfg.FieldMap.StackKey == "" {
		cfg.FieldMap.StackKey = defaultJSONFieldKeyStack
	}

	if cfg.FieldMap.MessageKey == "" {
		cfg.FieldMap.MessageKey = defaultJSONFieldKeyMessage
	}

	if cfg.DatetimeLayout == "" {
		cfg.DatetimeLayout = defaultDatetimeLayout
	}

	if cfg.TimestampFormat == 0 {
		cfg.TimestampFormat = defaultTimestampFormat
	}

	enc := new(EncoderJSON)
	enc.cfg = cfg

	return enc
}

// Encode encodes the given entry to the buffer.
func (enc *EncoderJSON) Encode(buf *buffer.Buffer, e Record) error { // nolint:funlen
	buf.WriteByte('{') // nolint:errcheck

	if e.Config.Datetime {
		buf.WriteString("\"")                         // nolint:errcheck
		buf.WriteString(enc.cfg.FieldMap.DatetimeKey) // nolint:errcheck
		buf.WriteString("\":\"")                      // nolint:errcheck
		buf.AppendTime(e.Time, enc.cfg.DatetimeLayout)
		buf.WriteString("\",") // nolint:errcheck
	}

	if e.Config.Timestamp {
		buf.WriteString("\"")                          // nolint:errcheck
		buf.WriteString(enc.cfg.FieldMap.TimestampKey) // nolint:errcheck
		buf.WriteString("\":\"")                       // nolint:errcheck
		switch enc.cfg.TimestampFormat {
		case TimestampFormatSeconds:
			buf.AppendInt(e.Time.Unix())
		case TimestampFormatNanoseconds:
			buf.AppendInt(e.Time.UnixNano())
		}
		buf.WriteString("\",") // nolint:errcheck
	}

	if levelStr := e.Level.String(); levelStr != "" {
		buf.WriteString("\"")                      // nolint:errcheck
		buf.WriteString(enc.cfg.FieldMap.LevelKey) // nolint:errcheck
		buf.WriteString("\":\"")                   // nolint:errcheck
		buf.WriteString(levelStr)                  // nolint:errcheck
		buf.WriteString("\",")                     // nolint:errcheck
	}

	if e.Config.Shortfile || e.Config.Longfile || e.Config.Stack {
		stackSkip := e.calldepth
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
				buf.WriteString("\"")                     // nolint:errcheck
				buf.WriteString(enc.cfg.FieldMap.FileKey) // nolint:errcheck
				buf.WriteString("\":\"")                  // nolint:errcheck
				buf.WriteString(file)                     // nolint:errcheck
				buf.WriteByte(':')                        // nolint:errcheck
				buf.AppendInt(int64(frames[0].Line))      // nolint:errcheck
				buf.WriteString("\",")                    // nolint:errcheck
			}
			if e.Config.Stack {
				buf.WriteString("\"")                      // nolint:errcheck
				buf.WriteString(enc.cfg.FieldMap.StackKey) // nolint:errcheck
				buf.WriteString("\":\"")                   // nolint:errcheck
				stackfmt := stack.NewStackFormatter(buf)
				stackfmt.FormatFrames(frames)
				buf.WriteString("\",") // nolint:errcheck
			}
		}
	}

	for _, field := range e.Fields {
		buf.WriteString("\"")      // nolint:errcheck
		buf.WriteString(field.Key) // nolint:errcheck
		buf.WriteString("\":\"")   // nolint:errcheck

		buf.WriteInterface(field.Val)
		buf.WriteString("\",") // nolint:errcheck
	}
	buf.WriteString("\"")                        // nolint:errcheck
	buf.WriteString(enc.cfg.FieldMap.MessageKey) // nolint:errcheck
	buf.WriteString("\":\"")                     // nolint:errcheck

	if needsQuote(e.Message) {
		buf.WriteString(strconv.Quote(e.Message))
	} else {
		buf.WriteString(e.Message) // nolint:errcheck
	}
	buf.WriteString("\"}") // nolint:errcheck
	buf.WriteNewLine()

	return nil
}

func needsQuote(s string) bool {
	for i := range s {
		c := s[i]
		if c < 0x20 || c > 0x7e || c == ' ' || c == '\\' || c == '"' || c == '\n' || c == '\r' {
			return true
		}
	}
	return false
}

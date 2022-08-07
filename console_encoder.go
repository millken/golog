package golog

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-json"

	"github.com/millken/golog/internal/buffer"
	"github.com/millken/golog/internal/stack"
)

var (
	_ Encoder = (*ConsoleEncoder)(nil)
)

// DefaultLineEnding is the default line ending used by the console encoder.
const DefaultLineEnding = '\n'

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	colorBold     = 1
	colorDarkGray = 90
)

const (
	consoleDefaultTimeFormat = time.RFC3339
)

// ConsoleEncoder encodes entries to the console.
type ConsoleEncoder struct {
	cfg ConsoleEncoderConfig
}

// NewConsoleEncoder returns a new console encoder.
func NewConsoleEncoder(cfg ConsoleEncoderConfig) *ConsoleEncoder {
	if len(cfg.PartsOrder) == 0 {
		cfg.PartsOrder = consoleDefaultPartsOrder()
	}
	return &ConsoleEncoder{
		cfg: cfg,
	}
}

// Encode encodes the entry and writes it to the writer.
func (o *ConsoleEncoder) Encode(e *Entry) ([]byte, error) {
	if e == nil {
		return nil, errors.New("nil entry")
	}
	var stacktraces string

	if o.cfg.DisableColor {
		e.SetFlag(FlagNoColor)
	}
	if e.HasFlag(FlagCaller) || e.HasFlag(FlagStacktrace) {
		stackSkip := defaultCallerSkip + e.CallerSkip() + o.cfg.CallerSkipFrame
		frames := stack.Tracer(stackSkip)

		if len(frames) > 0 {
			if e.HasFlag(FlagCaller) {
				frame := frames[0]
				e.SetCaller(frame.File + ":" + strconv.Itoa(frame.Line))
			}
			if e.HasFlag(FlagStacktrace) {
				buffer := buffer.Get()
				defer buffer.Free()
				stackfmt := stack.NewStackFormatter(buffer)
				stackfmt.FormatFrames(frames)
				stacktraces = buffer.String()
			}
		}
	}

	for _, p := range o.cfg.PartsOrder {
		if (p == CallerFieldName && !e.HasFlag(FlagCaller)) ||
			(p == ErrorStackFieldName && !e.HasFlag(FlagStacktrace) ||
				(p == TimestampFieldName && o.cfg.DisableTimestamp)) {
			continue
		}
		writePart(e, p)
		if p != o.cfg.PartsOrder[len(o.cfg.PartsOrder)-1] { // Skip space for last part
			e.WriteByte(' ')
		}
	}
	writeFields(e)
	if e.HasFlag(FlagStacktrace) {
		e.WriteByte(DefaultLineEnding)
		e.WriteString(stacktraces)
	}
	e.WriteByte(DefaultLineEnding)
	return e.Bytes(), nil
}

// writePart appends a formatted part to buf.
func writePart(e *Entry, p string) {
	switch p {
	case LevelFieldName:
		defaultFormatLevel(e)
	case TimestampFieldName:
		defaultFormatTimestamp(e, "")
	case MessageFieldName:
		defaultFormatMessage(e)
	case CallerFieldName:
		defaultFormatCaller(e)
	}
}

func writeFields(e *Entry) {
	if len(e.Fields) == 0 {
		return
	}
	for _, v := range e.Fields[:e.FieldsLength()] {
		e.WriteByte(' ')
		defaultFormatFieldName(e, v.Key)
		defaultFormatFieldValue(e, v.Val)
	}
}

func defaultFormatLevel(e *Entry) {
	var l string
	noColor := e.HasFlag(FlagNoColor)
	switch e.Level {
	case DEBUG:
		l = colorize("DBUG", colorCyan, noColor)
	case INFO:
		l = colorize("INFO", colorBlue, noColor)
	case WARNING:
		l = colorize("WARN", colorYellow, noColor)
	case ERROR:
		l = colorize("ERRO", colorRed, noColor)
	case FATAL:
		l = colorize(colorize("FATA", colorRed, noColor), colorBold, noColor)
	case PANIC:
		l = colorize(colorize("PNIC", colorDarkGray, noColor), colorBold, noColor)
	default:
		l = colorize("????", colorBold, noColor)
	}
	e.WriteString(l)
}

func defaultFormatTimestamp(e *Entry, timeFormat string) {
	if timeFormat == "" {
		timeFormat = consoleDefaultTimeFormat
	}
	e.Data = time.Now().AppendFormat(e.Data, timeFormat)
}

func defaultFormatMessage(e *Entry) {
	_, _ = e.WriteString(e.Message)
}

func defaultFormatCaller(e *Entry) {
	noColor := e.HasFlag(FlagNoColor)
	c := colorize(e.GetCaller(), colorBold, noColor)
	_, _ = e.WriteString(c)
}

func defaultFormatFieldName(e *Entry, name string) {
	noColor := e.HasFlag(FlagNoColor)
	_, _ = e.WriteString(colorize(name+"=", colorCyan, noColor))
}

func defaultFormatFieldValue(e *Entry, value interface{}) {
	switch fValue := value.(type) {
	case string:
		if needsQuote(fValue) {
			e.Data = append(e.Data, strconv.Quote(fValue)...)
		} else {
			e.Data = append(e.Data, fValue...)
		}
	case int:
		e.Data = strconv.AppendInt(e.Data, int64(fValue), 10)
	case int8:
		e.Data = strconv.AppendInt(e.Data, int64(fValue), 10)
	case int16:
		e.Data = strconv.AppendInt(e.Data, int64(fValue), 10)
	case int32:
		e.Data = strconv.AppendInt(e.Data, int64(fValue), 10)
	case int64:
		e.Data = strconv.AppendInt(e.Data, fValue, 10)
	case uint:
		e.Data = strconv.AppendUint(e.Data, uint64(fValue), 10)
	case uint8:
		e.Data = strconv.AppendUint(e.Data, uint64(fValue), 10)
	case uint16:
		e.Data = strconv.AppendUint(e.Data, uint64(fValue), 10)
	case uint32:
		e.Data = strconv.AppendUint(e.Data, uint64(fValue), 10)
	case uint64:
		e.Data = strconv.AppendUint(e.Data, fValue, 10)
	case float32:
		e.Data = strconv.AppendFloat(e.Data, float64(fValue), 'f', -1, 64)
	case float64:
		e.Data = strconv.AppendFloat(e.Data, float64(fValue), 'f', -1, 64)
	case bool:
		e.Data = strconv.AppendBool(e.Data, fValue)
	case error:
		e.Data = append(e.Data, fValue.Error()...)
	case []byte:
		e.Data = append(e.Data, fValue...)
	case time.Time:
		e.Data = fValue.AppendFormat(e.Data, consoleDefaultTimeFormat)
	case time.Duration:
		e.Data = append(e.Data, fValue.String()...)
	case json.Number:
		e.Data = append(e.Data, fValue.String()...)
	default:
		b, err := json.Marshal(fValue)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			e.Data = append(e.Data, b...)
		}
	}
}

// needsQuote returns true when the string s should be quoted in output.
func needsQuote(s string) bool {
	for i := range s {
		if s[i] < 0x20 || s[i] > 0x7e || s[i] == ' ' || s[i] == '\\' || s[i] == '"' {
			return true
		}
	}
	return false
}

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s string, c int, disabled bool) string {
	if disabled {
		return s
	}
	return "\x1b[" + strconv.Itoa(c) + "m" + s + "\x1b[0m"
}

func consoleDefaultPartsOrder() []string {
	return []string{
		TimestampFieldName,
		LevelFieldName,
		CallerFieldName,
		ErrorStackFieldName,
		MessageFieldName,
	}
}

package golog

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-json"

	"github.com/millken/golog/internal/buffer"
	"github.com/millken/golog/internal/fasttime"
	"github.com/millken/golog/internal/stack"
)

var (
	_ Encoder = (*TextEncoder)(nil)
)

// DefaultLineEnding is the default line ending used by the text encoder.
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
	textDefaultTimeFormat = time.RFC3339
)

// TextEncoder encodes entries to the text.
type TextEncoder struct {
	cfg TextEncoderConfig
}

// NewTextEncoder returns a new text encoder.
func NewTextEncoder(cfg TextEncoderConfig) *TextEncoder {
	if len(cfg.PartsOrder) == 0 {
		cfg.PartsOrder = textDefaultPartsOrder()
	}
	return &TextEncoder{
		cfg: cfg,
	}
}

// Encode encodes the entry and writes it to the writer.
func (o *TextEncoder) Encode(e *Entry) ([]byte, error) {
	if e == nil {
		return nil, errors.New("nil entry")
	}
	var stacktraces string

	if o.cfg.DisableColor {
		e.SetFlag(FlagNoColor)
	}
	if o.cfg.ShowModuleName {
		e.SetFlag(FlagName)
	}
	if e.HasFlag(FlagCaller) || e.HasFlag(FlagStacktrace) {
		stackSkip := DefaultCallerSkip + e.CallerSkip() + o.cfg.CallerSkipFrame
		frames := stack.Tracer(stackSkip, e.HasFlag(FlagStacktrace))

		if len(frames) > 0 {
			if e.HasFlag(FlagCaller) {
				frame := frames[0]
				e.SetCaller(frame.File + ":" + strconv.Itoa(frame.Line))
			}
			if e.HasFlag(FlagStacktrace) {
				buffer := buffer.Get()
				stackfmt := stack.NewStackFormatter(buffer)
				stackfmt.FormatFrames(frames)
				stacktraces = buffer.String()
				buffer.Free()
			}
		}
	}

	for _, p := range o.cfg.PartsOrder {
		if (p == CallerFieldName && !e.HasFlag(FlagCaller)) ||
			(p == ErrorStackFieldName && !e.HasFlag(FlagStacktrace) ||
				(p == TimestampFieldName && o.cfg.DisableTimestamp) ||
				(p == ModuleFieldName && !o.cfg.ShowModuleName)) {
			continue
		}
		writePart(e, p, &o.cfg)
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
func writePart(e *Entry, p string, cfg *TextEncoderConfig) {
	switch p {
	case LevelFieldName:
		defaultFormatLevel(e)
	case ModuleFieldName:
		defaultModuleName(e)
	case TimestampFieldName:
		timeFormat := cfg.TimeFormat
		if timeFormat == "" {
			timeFormat = textDefaultTimeFormat
		}
		defaultFormatTimestamp(e, timeFormat)
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
	noColor := e.HasFlag(FlagNoColor)
	switch e.Level {
	case DEBUG:
		ansiColorize("DBUG", colorCyan, noColor, e)
	case INFO:
		ansiColorize("INFO", colorBlue, noColor, e)
	case WARNING:
		ansiColorize("WARN", colorYellow, noColor, e)
	case ERROR:
		ansiColorize("ERRO", colorRed, noColor, e)
	case FATAL:
		ansiColorize("FATA", colorRed, noColor, e)
	case PANIC:
		ansiColorize("PNIC", colorDarkGray, noColor, e)
	default:
		ansiColorize("????", colorBold, noColor, e)
	}
	return
}

func defaultModuleName(e *Entry) {
	if e.HasFlag(FlagName) {
		_, _ = e.WriteString(e.Module)
	}
}

func defaultFormatTimestamp(e *Entry, timeFormat string) {
	e.Data = fasttime.Now().AppendFormat(e.Data, timeFormat)
}

func defaultFormatMessage(e *Entry) {
	_, _ = e.WriteString(e.Message)
}

func defaultFormatCaller(e *Entry) {
	if e.HasFlag(FlagNoColor) {
		_, _ = e.WriteString(e.GetCaller())
		return
	}
	ansiColorize(e.GetCaller(), colorBold, true, e)
	return
}

func defaultFormatFieldName(e *Entry, name string) {
	const equal string = "="
	if e.HasFlag(FlagNoColor) {
		e.Data = append(e.Data, name...)
		e.Data = append(e.Data, equal...)
		return
	}
	ansiColorize(name+equal, colorCyan, false, e)

	return
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
		e.Data = fValue.AppendFormat(e.Data, textDefaultTimeFormat)
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
		c := s[i]
		if c < 0x20 || c > 0x7e || c == ' ' || c == '\\' || c == '"' || c == '\n' || c == '\r' {
			return true
		}
	}
	return false
}

const (
	ansiReset = "\x1b[0m"
	ansiBold  = "\x1b["
)

func ansiColorize(s string, c int, disabled bool, e *Entry) {
	if disabled {
		_, _ = e.WriteString(s)
		return
	}
	_, _ = e.WriteString(ansiBold)
	_, _ = e.WriteString(strconv.Itoa(c))
	_, _ = e.WriteString("m")
	_, _ = e.WriteString(s)
	_, _ = e.WriteString(ansiReset)
	return
}

func textDefaultPartsOrder() []string {
	return []string{
		TimestampFieldName,
		LevelFieldName,
		ModuleFieldName,
		CallerFieldName,
		ErrorStackFieldName,
		MessageFieldName,
	}
}

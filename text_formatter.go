package golog

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/millken/golog/internal/buffer"
)

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

// TextFormatter formats log entries into text.
type TextFormatter struct {
	// NoColor disables the colorized output.
	NoColor bool
	// EnableCaller enabled caller
	EnableCaller bool
	// EnableStack enables stack trace
	EnableStack bool
	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp     bool
	CallerSkipFrameCount int

	// TimeFormat specifies the format for timestamp in output.
	TimeFormat string
	// PartsOrder defines the order of parts in output.
	PartsOrder []string

	// PartsExclude defines parts to not display in output.
	PartsExclude []string

	FormatTimestamp  func(*Entry)
	FormatLevel      func(*Entry)
	FormatCaller     func(*Entry)
	FormatMessage    func(*Entry)
	FormatFieldName  func(*Entry, string)
	FormatFieldValue func(*Entry, interface{})
}

// NewTextFormatter returns a new TextFormatter.
func NewTextFormatter() *TextFormatter {
	return &TextFormatter{
		TimeFormat: consoleDefaultTimeFormat,
		PartsOrder: consoleDefaultPartsOrder(),
	}
}

// Format renders a single log entry
func (f *TextFormatter) Format(entry *Entry) error {
	var stacktrace string
	stackDepth := stacktraceFirst
	if f.EnableStack {
		stackDepth = stacktraceFull
	}
	if f.EnableCaller || f.EnableStack {
		stack := captureStacktrace(entry.callerSkip+f.CallerSkipFrameCount, stackDepth)
		defer stack.Free()
		if stack.Count() > 0 {
			frame, more := stack.Next()
			if f.EnableCaller {
				c := frame.File + ":" + strconv.Itoa(frame.Line)
				entry.caller = c
			}
			if f.EnableStack {
				buffer := buffer.Get()
				defer buffer.Free()

				stackfmt := newStackFormatter(buffer)

				// We've already extracted the first frame, so format that
				// separately and defer to stackfmt for the rest.
				stackfmt.FormatFrame(frame)
				if more {
					stackfmt.FormatStack(stack)
				}
				stacktrace = buffer.String()
			}
		}
	}
	for _, p := range f.PartsOrder {
		f.writePart(entry, p)
	}
	f.writeFields(entry)
	if f.EnableStack {
		entry.WriteByte('\n')
		entry.WriteString(stacktrace)
	}
	return entry.WriteByte('\n')
}

// writeFields appends formatted key-fValueue pairs to buf.
func (f *TextFormatter) writeFields(entry *Entry) {
	if entry.FieldsLength() > 0 {
		_ = entry.WriteByte(' ')
	}

	i := 0
	for _, field := range entry.Fields[:entry.FieldsLength()] {
		name := field.Key
		fValue := field.Val
		i++
		if f.FormatFieldName == nil {
			f.defaultFormatFieldName(entry, name)
		} else {
			f.FormatFieldName(entry, name)
		}

		if f.FormatFieldValue == nil {
			f.defaultFormatFieldValue(entry, fValue)
		} else {
			f.FormatFieldValue(entry, fValue)
		}

		if i < entry.FieldsLength() { // Skip space for last field
			_ = entry.WriteByte(' ')
		}
	}
}

// writePart appends a formatted part to buf.
func (f *TextFormatter) writePart(entry *Entry, p string) {
	if f.DisableTimestamp && p == TimestampFieldName {
		return
	}
	if !f.EnableCaller && p == CallerFieldName {
		return
	}
	if !f.EnableStack && p == ErrorStackFieldName {
		return
	}

	if f.PartsExclude != nil && len(f.PartsExclude) > 0 {
		for _, exclude := range f.PartsExclude {
			if exclude == p {
				return
			}
		}
	}

	switch p {
	case LevelFieldName:
		if f.FormatLevel == nil {
			f.defaultFormatLevel(entry)
		} else {
			f.FormatLevel(entry)
		}
	case TimestampFieldName:
		if f.FormatTimestamp == nil {
			f.defaultFormatTimestamp(entry, f.TimeFormat)
		} else {
			f.FormatTimestamp(entry)
		}
	case MessageFieldName:
		if f.FormatMessage == nil {
			f.defaultFormatMessage(entry)
		} else {
			f.FormatMessage(entry)
		}
	case CallerFieldName:
		if f.FormatCaller == nil {
			f.defaultFormatCaller(entry)
		} else {
			f.FormatCaller(entry)
		}
	}
	if p != f.PartsOrder[len(f.PartsOrder)-1] { // Skip space for last part
		entry.WriteByte(' ')
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

func (f *TextFormatter) defaultFormatCaller(entry *Entry) {
	noColor := f.NoColor
	c := colorize(entry.caller, colorBold, noColor)
	_, _ = entry.WriteString(c)
}

func (f *TextFormatter) defaultFormatMessage(entry *Entry) {
	_, _ = entry.WriteString(entry.Message)
}

func (f *TextFormatter) defaultFormatFieldValue(entry *Entry, value interface{}) {
	switch fValue := value.(type) {
	case string:
		if needsQuote(fValue) {
			entry.Data = append(entry.Data, strconv.Quote(fValue)...)
		} else {
			entry.Data = append(entry.Data, fValue...)
		}
	case int:
		entry.Data = strconv.AppendInt(entry.Data, int64(fValue), 10)
	case int8:
		entry.Data = strconv.AppendInt(entry.Data, int64(fValue), 10)
	case int16:
		entry.Data = strconv.AppendInt(entry.Data, int64(fValue), 10)
	case int32:
		entry.Data = strconv.AppendInt(entry.Data, int64(fValue), 10)
	case int64:
		entry.Data = strconv.AppendInt(entry.Data, fValue, 10)
	case uint:
		entry.Data = strconv.AppendUint(entry.Data, uint64(fValue), 10)
	case uint8:
		entry.Data = strconv.AppendUint(entry.Data, uint64(fValue), 10)
	case uint16:
		entry.Data = strconv.AppendUint(entry.Data, uint64(fValue), 10)
	case uint32:
		entry.Data = strconv.AppendUint(entry.Data, uint64(fValue), 10)
	case uint64:
		entry.Data = strconv.AppendUint(entry.Data, fValue, 10)
	case float32:
		entry.Data = strconv.AppendFloat(entry.Data, float64(fValue), 'f', -1, 64)
	case float64:
		entry.Data = strconv.AppendFloat(entry.Data, float64(fValue), 'f', -1, 64)
	case bool:
		entry.Data = strconv.AppendBool(entry.Data, fValue)
	case error:
		entry.Data = append(entry.Data, fValue.Error()...)
	case []byte:
		entry.Data = append(entry.Data, fValue...)
	case time.Time:
		entry.Data = fValue.AppendFormat(entry.Data, consoleDefaultTimeFormat)
	case time.Duration:
		entry.Data = append(entry.Data, fValue.String()...)
	case json.Number:
		entry.Data = append(entry.Data, fValue.String()...)
	default:
		b, err := json.Marshal(fValue)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			entry.Data = append(entry.Data, b...)
		}
	}
}

func (f *TextFormatter) defaultFormatLevel(entry *Entry) {
	var l string
	ll := entry.Level
	noColor := f.NoColor
	switch ll {
	case Disabled:
		l = ""
	case DebugLevel:
		l = colorize("DBG", colorCyan, noColor)
	case InfoLevel:
		l = colorize("INF", colorBlue, noColor)
	case WarnLevel:
		l = colorize("WRN", colorYellow, noColor)
	case ErrorLevel:
		l = colorize("ERR", colorRed, noColor)
	case FatalLevel:
		l = colorize(colorize("FTL", colorRed, noColor), colorBold, noColor)
	case PanicLevel:
		l = colorize(colorize("PNC", colorDarkGray, noColor), colorBold, noColor)
	default:
		l = colorize("?????", colorBold, noColor)
	}
	entry.WriteString(l)
}

func (f *TextFormatter) defaultFormatFieldName(entry *Entry, name string) {
	entry.Data = append(entry.Data, colorize(name+"=", colorCyan, f.NoColor)...)
}

func (f *TextFormatter) defaultFormatTimestamp(entry *Entry, timeFormat string) {
	if timeFormat == "" {
		timeFormat = consoleDefaultTimeFormat
	}
	entry.Data = entry.Timestamp.AppendFormat(entry.Data, timeFormat)
}

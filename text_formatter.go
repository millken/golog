package golog

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

type TextFormatter struct {
	// NoColor disables the colorized output.
	NoColor bool
	// EnableCaller enabled caller
	EnableCaller         bool
	CallerSkipFrameCount int

	// TimeFormat specifies the format for timestamp in output.
	TimeFormat string

	// PartsOrder defines the order of parts in output.
	PartsOrder []string

	// PartsExclude defines parts to not display in output.
	PartsExclude []string

	FormatTimestamp     func(*Entry)
	FormatLevel         func(*Entry)
	FormatCaller        func(*Entry)
	FormatMessage       func(*Entry)
	FormatFieldName     func(*Entry, string)
	FormatFieldValue    func(*Entry, interface{})
	FormatErrFieldName  Formatter
	FormatErrFieldValue Formatter
}

func (f *TextFormatter) Format(entry *Entry) error {
	if f.PartsOrder == nil {
		f.PartsOrder = consoleDefaultPartsOrder()
	}

	for _, p := range f.PartsOrder {
		f.writePart(entry, p)
	}

	f.writeFields(entry)

	return entry.WriteByte('\n')
}

// writeFields appends formatted key-fValueue pairs to buf.
func (f *TextFormatter) writeFields(entry *Entry) {
	if cap(entry.Fields) > 0 {
		entry.WriteByte(' ')
	}

	i := 0
	for _, field := range entry.Fields {
		name := b2s(field.key)
		fValue := field.value
		i++
		if name == ErrorFieldName {
			// if f.FormatErrFieldName == nil {
			// 	consoleDefaultFormatErrFieldName(f.NoColor)
			// } else {
			// 	fn = f.FormatErrFieldName
			// }

			// if f.FormatErrFieldValue == nil {
			// 	fv = consoleDefaultFormatErrFieldValue(f.NoColor)
			// } else {
			// 	fv = f.FormatErrFieldValue
			// }
		} else {
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
		}

		if i < len(entry.Fields) { // Skip space for last field
			entry.WriteByte(' ')
		}
	}
}

// writePart appends a formatted part to buf.
func (f *TextFormatter) writePart(entry *Entry, p string) {

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
		if f.EnableCaller {
			if f.FormatCaller == nil {
				f.defaultFormatCaller(entry)
			} else {
				f.FormatCaller(entry)
			}
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
		MessageFieldName,
	}
}

func (f *TextFormatter) defaultFormatCaller(entry *Entry) {
	var c string
	skip := f.CallerSkipFrameCount
	if skip == 0 {
		skip = CallerSkipFrameCount + 2
	}
	noColor := f.NoColor
	file, line := entry.GetCaller(skip)
	c = file + ":" + strconv.Itoa(line)
	if len(c) > 0 {
		if cwd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(cwd, c); err == nil {
				c = rel
			}
		}
		c = colorize(c, colorBold, noColor) + colorize(" >", colorCyan, noColor)
	}
	entry.WriteString(c)
}

func (f *TextFormatter) defaultFormatMessage(entry *Entry) {
	entry.Write(entry.Data)
}

func (f *TextFormatter) defaultFormatFieldValue(entry *Entry, value interface{}) {
	switch fValue := value.(type) {
	case string:
		if needsQuote(fValue) {
			entry.Formatted = append(entry.Formatted, strconv.Quote(fValue)...)
		} else {
			entry.Formatted = append(entry.Formatted, fValue...)
		}
	case int:
		entry.Formatted = strconv.AppendInt(entry.Formatted, int64(fValue), 10)
	case int8:
		entry.Formatted = strconv.AppendInt(entry.Formatted, int64(fValue), 10)
	case int16:
		entry.Formatted = strconv.AppendInt(entry.Formatted, int64(fValue), 10)
	case int32:
		entry.Formatted = strconv.AppendInt(entry.Formatted, int64(fValue), 10)
	case int64:
		entry.Formatted = strconv.AppendInt(entry.Formatted, fValue, 10)
	case uint:
		entry.Formatted = strconv.AppendUint(entry.Formatted, uint64(fValue), 10)
	case uint8:
		entry.Formatted = strconv.AppendUint(entry.Formatted, uint64(fValue), 10)
	case uint16:
		entry.Formatted = strconv.AppendUint(entry.Formatted, uint64(fValue), 10)
	case uint32:
		entry.Formatted = strconv.AppendUint(entry.Formatted, uint64(fValue), 10)
	case uint64:
		entry.Formatted = strconv.AppendUint(entry.Formatted, fValue, 10)
	case float32:
		entry.Formatted = strconv.AppendFloat(entry.Formatted, float64(fValue), 'f', -1, 64)
	case float64:
		entry.Formatted = strconv.AppendFloat(entry.Formatted, float64(fValue), 'f', -1, 64)
	case bool:
		entry.Formatted = strconv.AppendBool(entry.Formatted, fValue)
	case error:
		entry.Formatted = append(entry.Formatted, fValue.Error()...)
	case []byte:
		entry.Formatted = append(entry.Formatted, fValue...)
	case time.Time:
		entry.Formatted = fValue.AppendFormat(entry.Formatted, consoleDefaultTimeFormat)
	case json.Number:
		entry.Formatted = append(entry.Formatted, fValue.String()...)
	default:
		b, err := json.Marshal(fValue)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			entry.Formatted = append(entry.Formatted, b...)
		}
	}
}

func (f *TextFormatter) defaultFormatLevel(entry *Entry) {
	var l string
	ll := entry.Level
	noColor := f.NoColor
	switch ll {
	case TraceLevel:
		l = colorize(ll.String(), colorMagenta, noColor)
	case DebugLevel:
		l = colorize(ll.String(), colorYellow, noColor)
	case InfoLevel:
		l = colorize(ll.String(), colorGreen, noColor)
	case WarnLevel:
		l = colorize(ll.String(), colorRed, noColor)
	case ErrorLevel:
		l = colorize(colorize(ll.String(), colorRed, noColor), colorBold, noColor)
	case FatalLevel:
		l = colorize(colorize(ll.String(), colorRed, noColor), colorBold, noColor)
	case PanicLevel:
		l = colorize(colorize(ll.String(), colorRed, noColor), colorBold, noColor)
	default:
		l = colorize("???", colorBold, noColor)
	}
	entry.WriteString(l)
}

func (f *TextFormatter) defaultFormatFieldName(entry *Entry, name string) {
	entry.Formatted = append(entry.Formatted, colorize(name+"=", colorCyan, f.NoColor)...)
}

func (f *TextFormatter) defaultFormatTimestamp(entry *Entry, timeFormat string) {
	if timeFormat == "" {
		timeFormat = consoleDefaultTimeFormat
	}
	entry.Formatted = entry.Timestamp.AppendFormat(entry.Formatted, timeFormat)
}

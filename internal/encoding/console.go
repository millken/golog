package encoding

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/millken/golog/internal/config"
	"github.com/millken/golog/internal/log"
	"github.com/millken/golog/internal/stacktrace"
)

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

var (
	defaultSkip int = 5
)

type console struct {
	cfg config.ConsoleConfig
}

// NewConsoleEncoder returns a new console encoder.
func NewConsole(cfg config.ConsoleConfig) *console {
	if len(cfg.PartsOrder) == 0 {
		cfg.PartsOrder = consoleDefaultPartsOrder()
	}
	return &console{
		cfg: cfg,
	}
}

// Enocde encodes the entry and writes it to the writer.
func (c *console) Encode(e *log.Entry) ([]byte, error) {
	if e == nil {
		return nil, errors.New("nil entry")
	}
	//var stacktraces string
	stackDepth := stacktrace.StacktraceFirst

	if e.HasFlag(log.FlagCaller) {
		stack := stacktrace.Capture(defaultSkip, stackDepth)
		defer stack.Free()
		if stack.Count() > 0 {
			frame, _ := stack.Next()
			if e.HasFlag(log.FlagCaller) {
				c := frame.File + ":" + strconv.Itoa(frame.Line)
				e.SetCaller(c)
			}
		}
	}

	for _, p := range c.cfg.PartsOrder {
		writePart(e, p)
		if p != c.cfg.PartsOrder[len(c.cfg.PartsOrder)-1] { // Skip space for last part
			e.WriteByte(' ')
		}
	}
	writeFields(e)
	e.WriteByte(DefaultLineEnding)
	return e.Bytes(), nil
}

// writePart appends a formatted part to buf.
func writePart(e *log.Entry, p string) {
	switch p {
	case log.LevelFieldName:
		defaultFormatLevel(e)
	case log.TimestampFieldName:
		defaultFormatTimestamp(e, "")
	case log.MessageFieldName:
		defaultFormatMessage(e)
	case log.CallerFieldName:
		defaultFormatCaller(e)
	}
}

func writeFields(e *log.Entry) {
	if len(e.Fields) == 0 {
		return
	}
	for _, v := range e.Fields[:e.FieldsLength()] {
		e.WriteByte(' ')
		defaultFormatFieldName(e, v.Key)
		defaultFormatFieldValue(e, v)
	}
	e.WriteByte(DefaultLineEnding)
}

func defaultFormatLevel(e *log.Entry) {
	var l string
	noColor := !e.HasFlag(log.FlagColor)
	switch e.Level {
	case log.DEBUG:
		l = colorize("DBG", colorCyan, noColor)
	case log.INFO:
		l = colorize("INF", colorBlue, noColor)
	case log.WARNING:
		l = colorize("WRN", colorYellow, noColor)
	case log.ERROR:
		l = colorize("ERR", colorRed, noColor)
	case log.FATAL:
		l = colorize(colorize("FTL", colorRed, noColor), colorBold, noColor)
	case log.PANIC:
		l = colorize(colorize("PNC", colorDarkGray, noColor), colorBold, noColor)
	default:
		l = colorize("???", colorBold, noColor)
	}
	e.WriteString(l)
}

func defaultFormatTimestamp(e *log.Entry, timeFormat string) {
	if timeFormat == "" {
		timeFormat = consoleDefaultTimeFormat
	}
	e.Data = time.Now().AppendFormat(e.Data, timeFormat)
}

func defaultFormatMessage(e *log.Entry) {
	_, _ = e.WriteString(e.Message)
}

func defaultFormatCaller(e *log.Entry) {
	noColor := !e.HasFlag(log.FlagColor)
	c := colorize(e.GetCaller(), colorBold, noColor)
	_, _ = e.WriteString(c)
}

func defaultFormatFieldName(e *log.Entry, name string) {
	noColor := !e.HasFlag(log.FlagColor)
	_, _ = e.WriteString(colorize(name+"=", colorCyan, noColor))
}

func defaultFormatFieldValue(e *log.Entry, value interface{}) {
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
		log.TimestampFieldName,
		log.LevelFieldName,
		log.CallerFieldName,
		log.ErrorStackFieldName,
		log.MessageFieldName,
	}
}

// // getCallerInfo going through runtime caller frames to determine the caller of logger function by filtering
// // internal logging library functions.
// func getCallerInfo(level log.Level) string {

// 	const (
// 		// search MAXCALLERS caller frames for the real caller,
// 		// MAXCALLERS defines maximum number of caller frames needed to be recorded to find the actual caller frame
// 		MAXCALLERS = 6
// 		// skip SKIPCALLERS frames when determining the real caller
// 		// SKIPCALLERS is the number of stack frames to skip before recording caller frames,
// 		// this is mainly used to filter logger library functions in caller frames
// 		SKIPCALLERS      = 5
// 		NOTFOUND         = "n/a"
// 		DEFAULTLOGPREFIX = "log.(*Log)"
// 	)

// 	fpcs := make([]uintptr, MAXCALLERS)

// 	n := runtime.Callers(SKIPCALLERS, fpcs)
// 	if n == 0 {
// 		return fmt.Sprintf(callerInfoFormatter, NOTFOUND)
// 	}

// 	frames := runtime.CallersFrames(fpcs[:n])
// 	loggerFrameFound := false

// 	for f, more := frames.Next(); more; f, more = frames.Next() {
// 		_, fnName := filepath.Split(f.Function)

// 		if f.Func == nil || f.Function == "" {
// 			fnName = NOTFOUND // not a function or unknown
// 		}

// 		if loggerFrameFound {
// 			return fmt.Sprintf(callerInfoFormatter, fnName)
// 		}

// 		if strings.HasPrefix(fnName, DEFAULTLOGPREFIX) {
// 			loggerFrameFound = true

// 			continue
// 		}

// 		return fmt.Sprintf(callerInfoFormatter, fnName)
// 	}

// 	return fmt.Sprintf(callerInfoFormatter, NOTFOUND)
// }

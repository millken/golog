package golog

import (
	"strconv"

	"github.com/millken/golog/internal/buffer"
)

// JSONFormatter is a formatter that outputs JSON-encoded log messages.
type JSONFormatter struct {
	// EnableCaller enabled caller
	EnableCaller bool
	// EnableStack enables stack trace
	EnableStack          bool
	CallerSkipFrameCount int
	DisableTimestamp     bool
}

// Format formats the log entry.
func (f *JSONFormatter) Format(entry *Entry) error {
	entry.Data = enc.AppendBeginMarker(entry.Data)
	if !f.DisableTimestamp {
		entry.Data = appendKeyVal(entry.Data, TimestampFieldName, &entry.Timestamp)
	}
	entry.Data = appendKeyVal(entry.Data, LevelFieldName, entry.Level.String())
	entry.Data = appendKeyVal(entry.Data, MessageFieldName, &entry.Message)

	stackDepth := stacktraceFirst
	if f.EnableStack {
		stackDepth = stacktraceFull
	}
	if f.EnableCaller || f.EnableStack {
		stack := captureStacktrace(entry.callerSkip, stackDepth)
		defer stack.Free()
		if stack.Count() > 0 {
			frame, more := stack.Next()
			if f.EnableCaller {
				c := frame.File + ":" + strconv.Itoa(frame.Line)
				entry.Data = appendKeyVal(entry.Data, CallerFieldName, c)
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
				entry.Data = appendKeyVal(entry.Data, ErrorStackFieldName, buffer.String())
			}
		}
	}
	entry.Data = appendFields(entry.Data, entry.Fields[:entry.FieldsLength()]...)
	entry.Data = enc.AppendEndMarker(entry.Data)
	entry.Data = enc.AppendLineBreak(entry.Data)
	return nil
}

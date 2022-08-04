package encoding

import (
	"errors"
	"strconv"
	"time"

	"github.com/millken/golog/internal/config"
	"github.com/millken/golog/internal/log"
	"github.com/millken/golog/internal/stacktrace"
)

var (
	_ log.Encoder = (*JSONEncoder)(nil)
)

// JSONEncoder encodes entries as JSON.
type JSONEncoder struct {
	cfg config.JSONEncoderConfig
}

// NewJSONEncoder returns a new JSONEncoder.
func NewJSONEncoder(cfg config.JSONEncoderConfig) *JSONEncoder {
	return &JSONEncoder{
		cfg: cfg,
	}
}

// Encode encodes the entry and writes it to the writer.
func (o *JSONEncoder) Encode(e *log.Entry) ([]byte, error) {
	if e == nil {
		return nil, errors.New("nil entry")
	}
	e.Data = enc.AppendBeginMarker(e.Data)
	if !o.cfg.DisableTimestamp {
		e.Data = appendKeyVal(e.Data, log.TimestampFieldName, time.Now())
	}
	e.Data = appendKeyVal(e.Data, log.LevelFieldName, e.Level.String())
	e.Data = appendKeyVal(e.Data, log.MessageFieldName, &e.Message)

	//var stacktraces string
	stackDepth := stacktrace.StacktraceFirst

	if e.HasFlag(log.FlagCaller) {
		stack := stacktrace.Capture(defaultSkip, stackDepth)
		defer stack.Free()
		if stack.Count() > 0 {
			frame, _ := stack.Next()
			if e.HasFlag(log.FlagCaller) {
				c := frame.File + ":" + strconv.Itoa(frame.Line)
				e.Data = appendKeyVal(e.Data, log.CallerFieldName, c)
			}
		}
	}

	e.Data = appendFields(e.Data, e.Fields[:e.FieldsLength()]...)
	e.Data = enc.AppendEndMarker(e.Data)
	e.Data = enc.AppendLineBreak(e.Data)
	return e.Bytes(), nil
}

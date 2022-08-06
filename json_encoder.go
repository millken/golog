package golog

import (
	"errors"
	"runtime"
	"strconv"
	"time"

	"github.com/millken/golog/internal/buffer"
	"github.com/millken/golog/internal/stack"
)

var (
	_ Encoder = (*JSONEncoder)(nil)
)

// JSONEncoder encodes entries as JSON.
type JSONEncoder struct {
	cfg JSONEncoderConfig
}

// NewJSONEncoder returns a new JSONEncoder.
func NewJSONEncoder(cfg JSONEncoderConfig) *JSONEncoder {
	return &JSONEncoder{
		cfg: cfg,
	}
}

// Encode encodes the entry and writes it to the writer.
func (o *JSONEncoder) Encode(e *Entry) ([]byte, error) {
	if e == nil {
		return nil, errors.New("nil entry")
	}
	e.Data = enc.AppendBeginMarker(e.Data)
	if !o.cfg.DisableTimestamp {
		e.Data = appendKeyVal(e.Data, TimestampFieldName, time.Now())
	}
	e.Data = appendKeyVal(e.Data, LevelFieldName, e.Level.String())
	e.Data = appendKeyVal(e.Data, MessageFieldName, &e.Message)

	var frames []runtime.Frame
	if e.HasFlag(FlagCaller) || e.HasFlag(FlagStacktrace) {
		stackSkip := defaultCallerSkip + e.CallerSkip() + o.cfg.CallerSkipFrame
		frames = stack.Tracer(stackSkip)
	}

	if len(frames) > 0 {
		if e.HasFlag(FlagCaller) {
			frame := frames[0]
			c := frame.File + ":" + strconv.Itoa(frame.Line)
			e.Data = appendKeyVal(e.Data, CallerFieldName, c)
		}
		if e.HasFlag(FlagStacktrace) {
			buffer := buffer.Get()
			defer buffer.Free()
			stackfmt := stack.NewStackFormatter(buffer)
			stackfmt.FormatFrames(frames)
			e.Data = appendKeyVal(e.Data, ErrorStackFieldName, buffer.String())
		}
	}

	e.Data = appendFields(e.Data, e.Fields[:e.FieldsLength()]...)
	e.Data = enc.AppendEndMarker(e.Data)
	e.Data = enc.AppendLineBreak(e.Data)
	return e.Bytes(), nil
}

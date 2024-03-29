package golog

import (
	"errors"
	"runtime"
	"strconv"

	"github.com/millken/golog/internal/buffer"
	"github.com/millken/golog/internal/fasttime"
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
		e.Data = enc.AppendKey(e.Data, TimestampFieldName)
		e.Data = enc.AppendTime(e.Data, fasttime.Now(), TimeFieldFormat)
	}
	e.Data = enc.AppendKey(e.Data, LevelFieldName)
	e.Data = enc.AppendString(e.Data, e.Level.String())
	if o.cfg.ShowModuleName {
		e.Data = enc.AppendKey(e.Data, ModuleFieldName)
		e.Data = enc.AppendString(e.Data, e.Module)
	}
	e.Data = enc.AppendKey(e.Data, MessageFieldName)
	e.Data = enc.AppendString(e.Data, e.Message)

	var frames []runtime.Frame
	if e.HasFlag(FlagCaller) || e.HasFlag(FlagStacktrace) {
		stackSkip := DefaultCallerSkip + e.CallerSkip() + o.cfg.CallerSkipFrame
		frames = stack.Tracer(stackSkip, e.HasFlag(FlagStacktrace))
	}

	if len(frames) > 0 {
		if e.HasFlag(FlagCaller) {
			frame := frames[0]
			e.Data = enc.AppendKey(e.Data, CallerFieldName)
			e.Data = enc.AppendString(e.Data, frame.File+":"+strconv.Itoa(frame.Line))
		}
		if e.HasFlag(FlagStacktrace) {
			buffer := buffer.Get()
			stackfmt := stack.NewStackFormatter(buffer)
			stackfmt.FormatFrames(frames)
			e.Data = enc.AppendKey(e.Data, ErrorStackFieldName)
			e.Data = enc.AppendBytes(e.Data, buffer.Bytes())
			buffer.Free()
		}
	}
	for _, field := range e.Fields[:e.FieldsLength()] {
		e.Data = enc.AppendKey(e.Data, field.Key)
		e.Data = appendVal(e.Data, field.Val)
	}
	e.Data = enc.AppendEndMarker(e.Data)
	e.Data = enc.AppendLineBreak(e.Data)
	return e.Bytes(), nil
}

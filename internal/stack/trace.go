package stack

import (
	"runtime"

	"github.com/millken/golog/internal/buffer"
)

var (
	// MaxCallers is the maximum number of callers to include in the stack.
	MaxCallers = 20
)

// Tracer returns a slice of Frames, calling runtime.Callers.
func Tracer(skip int) []runtime.Frame {
	var stack []runtime.Frame

	fpcs := make([]uintptr, MaxCallers)

	//+2 to skip Tracer and runtime.Callers.
	n := runtime.Callers(skip+2, fpcs)
	if n == 0 {
		return nil
	}

	frames := runtime.CallersFrames(fpcs[:n])
	for f, more := frames.Next(); more; f, more = frames.Next() {
		stack = append(stack, f)
	}
	if len(stack) == 0 {
		return nil
	}

	return stack
}

// stackFormatter formats a stack trace into a readable string representation.
type stackFormatter struct {
	b        *buffer.Buffer
	nonEmpty bool // whehther we've written at least one frame already
}

// NewStackFormatter builds a new stackFormatter.
func NewStackFormatter(b *buffer.Buffer) stackFormatter {
	return stackFormatter{b: b}
}

// FormatFrames formats all remaining frames in the provided frames -- minus
// the final runtime.main/runtime.goexit frame.
func (sf *stackFormatter) FormatFrames(frames []runtime.Frame) {
	for _, f := range frames {
		sf.FormatFrame(f)
	}
}

// FormatFrame formats the given frame.
func (sf *stackFormatter) FormatFrame(frame runtime.Frame) {
	if sf.nonEmpty {
		sf.b.AppendByte('\n')
	}
	sf.nonEmpty = true
	sf.b.AppendString(frame.Function)
	sf.b.AppendByte('\n')
	sf.b.AppendByte('\t')
	sf.b.AppendString(frame.File)
	sf.b.AppendByte(':')
	sf.b.AppendInt(int64(frame.Line))
}
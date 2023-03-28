package stack

import (
	"runtime"
	"sync"

	"github.com/millken/golog/internal/buffer"
)

const (
	offset32 = 2166136261
	prime32  = 16777619
)

var (
	m sync.Map
)

// Tracer returns a slice of Frames, calling runtime.Callers.
func Tracer(skip int, stacktrace bool) []runtime.Frame {
	const depth = 20
	var stack []runtime.Frame
	if !stacktrace {
		pc, file, line, ok := runtime.Caller(skip + 1)
		if !ok {
			return nil
		}
		stack = append(stack, runtime.Frame{
			PC:   pc,
			File: file,
			Line: line,
		})
		return stack
	}

	//the maximum number of callers to include in the stack.
	var fpcs [depth]uintptr

	//+2 to skip Tracer and runtime.Callers.
	n := runtime.Callers(skip+2, fpcs[:])
	if n == 0 {
		return nil
	}
	hash := uint(offset32)
	for _, pc := range fpcs[:n] {
		hash ^= uint(pc)
		hash *= prime32
	}
	if item, ok := m.Load(hash); ok {
		return item.([]runtime.Frame)
	}
	frames := runtime.CallersFrames(fpcs[:n])
	for f, more := frames.Next(); more; f, more = frames.Next() {
		stack = append(stack, f)
	}
	if len(stack) == 0 {
		return nil
	}
	m.Store(hash, stack)
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

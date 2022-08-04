package stacktrace

import (
	"runtime"
	"sync"

	"github.com/millken/golog/internal/buffer"
)

var _stacktracePool = sync.Pool{
	New: func() interface{} {
		return &StackTrace{
			storage: make([]uintptr, 64),
		}
	},
}

// StackTrace represents
type StackTrace struct {
	pcs    []uintptr // program counters; always a subslice of storage
	frames *runtime.Frames

	// The size of pcs varies depending on requirements:
	// it will be one if the only the first frame was requested,
	// and otherwise it will reflect the depth of the call stack.
	//
	// storage decouples the slice we need (pcs) from the slice we pool.
	// We will always allocate a reasonably large storage, but we'll use
	// only as much of it as we need.
	storage []uintptr
}

// stacktraceDepth specifies how deep of a stack trace should be captured.
type stacktraceDepth int

const (
	// StacktraceFirst captures only the first frame.
	StacktraceFirst stacktraceDepth = iota

	// StacktraceFull captures the entire call stack, allocating more
	// storage for it if needed.
	StacktraceFull
)

// Capture captures a stack trace of the specified depth, skipping
// the provided number of frames. skip=0 identifies the caller of
// captureStacktrace.
//
// The caller must call Free on the returned stacktrace after using it.
func Capture(skip int, depth stacktraceDepth) *StackTrace {
	stack := _stacktracePool.Get().(*StackTrace)

	switch depth {
	case StacktraceFirst:
		stack.pcs = stack.storage[:1]
	case StacktraceFull:
		stack.pcs = stack.storage
	}

	// Unlike other "skip"-based APIs, skip=0 identifies runtime.Callers
	// itself. +2 to skip captureStacktrace and runtime.Callers.
	numFrames := runtime.Callers(
		skip+2,
		stack.pcs,
	)

	// runtime.Callers truncates the recorded stacktrace if there is no
	// room in the provided slice. For the full stack trace, keep expanding
	// storage until there are fewer frames than there is room.
	if depth == StacktraceFull {
		pcs := stack.pcs
		for numFrames == len(pcs) {
			pcs = make([]uintptr, len(pcs)*2)
			numFrames = runtime.Callers(skip+2, pcs)
		}

		// Discard old storage instead of returning it to the pool.
		// This will adjust the pool size over time if stack traces are
		// consistently very deep.
		stack.storage = pcs
		stack.pcs = pcs[:numFrames]
	} else {
		stack.pcs = stack.pcs[:numFrames]
	}

	stack.frames = runtime.CallersFrames(stack.pcs)
	return stack
}

// Free releases resources associated with this stacktrace
// and returns it back to the pool.
func (st *StackTrace) Free() {
	st.frames = nil
	st.pcs = nil
	_stacktracePool.Put(st)
}

// Count reports the total number of frames in this stacktrace.
// Count DOES NOT change as Next is called.
func (st *StackTrace) Count() int {
	return len(st.pcs)
}

// Next returns the next frame in the stack trace,
// and a boolean indicating whether there are more after it.
func (st *StackTrace) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next()
}

// TakeStacktrace takes a stacktrace of the specified depth, skipping the
func TakeStacktrace(skip int) string {
	stack := Capture(skip+1, StacktraceFull)
	defer stack.Free()

	buffer := buffer.Get()
	defer buffer.Free()

	stackfmt := NewStackFormatter(buffer)
	stackfmt.FormatStack(stack)
	return buffer.String()
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

// FormatStack formats all remaining frames in the provided stacktrace -- minus
// the final runtime.main/runtime.goexit frame.
func (sf *stackFormatter) FormatStack(stack *StackTrace) {
	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	for frame, more := stack.Next(); more; frame, more = stack.Next() {
		sf.FormatFrame(frame)
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

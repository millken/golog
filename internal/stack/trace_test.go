package stack

import (
	"testing"

	"github.com/millken/golog/internal/buffer"
	"github.com/stretchr/testify/require"
)

func TestTrace(t *testing.T) {
	require := require.New(t)
	frames := Tracer(0, true)
	if len(frames) == 0 {
		t.Fatal("no frames")
	}
	frame := frames[0]
	require.NotNil(frame)
	require.Contains(frame.File, "stack/trace_test.go")

	b := buffer.Get()
	defer b.Free()
	sf := NewStackFormatter(b)
	sf.FormatFrames(frames)
	require.Contains(b.String(), "stack/trace_test.go")
}

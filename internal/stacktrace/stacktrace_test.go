package stacktrace

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTakeStacktrace(t *testing.T) {
	trace := TakeStacktrace(0)
	lines := strings.Split(trace, "\n")
	require.NotEmpty(t, lines, "Expected stacktrace to have at least one frame.")
	assert.Contains(
		t,
		lines[0],
		"stacktrace.TestTakeStacktrace",
		"Expected stacktrace to start with the test.",
	)
}

func TestTakeStacktraceWithSkip(t *testing.T) {
	trace := TakeStacktrace(1)
	lines := strings.Split(trace, "\n")
	require.NotEmpty(t, lines, "Expected stacktrace to have at least one frame.")
	assert.Contains(
		t,
		lines[0],
		"testing.",
		"Expected stacktrace to start with the test runner (skipping our own frame).",
	)
}

func TestTakeStacktraceWithSkipInnerFunc(t *testing.T) {
	var trace string
	func() {
		trace = TakeStacktrace(2)
	}()
	lines := strings.Split(trace, "\n")
	require.NotEmpty(t, lines, "Expected stacktrace to have at least one frame.")
	assert.Contains(
		t,
		lines[0],
		"testing.",
		"Expected stacktrace to start with the test function (skipping the test function).",
	)
}

func TestTakeStacktraceDeepStack(t *testing.T) {
	const (
		N                  = 500
		withStackDepthName = "stacktrace.withStackDepth"
	)
	withStackDepth(N, func() {
		trace := TakeStacktrace(0)
		for found := 0; found < N; found++ {
			i := strings.Index(trace, withStackDepthName)
			if i < 0 {
				t.Fatalf(`expected %v occurrences of %q, found %d`,
					N, withStackDepthName, found)
			}
			trace = trace[i+len(withStackDepthName):]
		}
	})
}

func BenchmarkTakeStacktrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TakeStacktrace(0)
	}
}

func withStackDepth(depth int, f func()) {
	var recurse func(rune) rune
	recurse = func(r rune) rune {
		if r > 0 {
			bytes.Map(recurse, []byte(string([]rune{r - 1})))
		} else {
			f()
		}
		return 0
	}
	recurse(rune(depth))
}

package golog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConsole(t *testing.T) {

	module := "test-module"
	cs := NewConsoleEncoder(ConsoleEncoderConfig{})
	_, err := cs.Encode(nil)
	require.Error(t, err)
	e := &Entry{
		Module:  module,
		Level:   INFO,
		Message: "test",
	}
	e.SetFlag(FlagCaller)
	defaultCallerSkip = 1
	b, err := cs.Encode(e)
	require.NoError(t, err)
	require.Contains(t, string(b), "test")
	require.Contains(t, string(b), "INF")
	require.Contains(t, string(b), "console_encoder_test.go")

}

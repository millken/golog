package encoding

import (
	"testing"

	"github.com/millken/golog/config"
	"github.com/millken/golog/log"
	"github.com/stretchr/testify/require"
)

func TestConsole(t *testing.T) {

	module := "test-module"
	cs := NewConsoleEncoder(config.ConsoleEncoderConfig{})
	_, err := cs.Encode(nil)
	require.Error(t, err)
	e := &log.Entry{
		Module:  module,
		Level:   log.INFO,
		Message: "test",
	}
	e.SetFlag(log.FlagCaller)
	defaultCallerSkip = 1
	b, err := cs.Encode(e)
	require.NoError(t, err)
	require.Contains(t, string(b), "test")
	require.Contains(t, string(b), "INF")
	require.Contains(t, string(b), "console_encoder_test.go")

}

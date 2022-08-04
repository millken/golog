package encoding

import (
	"testing"

	"github.com/millken/golog/internal/config"
	"github.com/millken/golog/internal/log"
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
	defaultSkip = 1
	b, err := cs.Encode(e)
	require.NoError(t, err)
	require.Contains(t, string(b), "test")
	require.Contains(t, string(b), "INF")
	require.Contains(t, string(b), "console_test.go")

}
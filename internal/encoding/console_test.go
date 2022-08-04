package encoding

import (
	"testing"
	"time"

	"github.com/millken/golog/internal/config"
	"github.com/millken/golog/internal/log"
	"github.com/millken/golog/internal/meta"
	"github.com/stretchr/testify/require"
)

func TestConsole(t *testing.T) {

	module := "test-module"
	level := log.INFO
	meta.ShowCallerInfo(module, level)
	cs := NewConsole(config.ConsoleConfig{})
	_, err := cs.Encode(nil)
	require.Error(t, err)
	e := &log.Entry{
		Timestamp: time.Now(),
		Module:    module,
		Level:     log.INFO,
		Message:   "test",
	}
	b, err := cs.Encode(e)
	require.NoError(t, err)
	require.Contains(t, string(b), "test")
	require.Contains(t, string(b), "INF")
	require.Contains(t, string(b), "console_test.go")

}

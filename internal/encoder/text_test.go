package encoder

import (
	"testing"
	"time"

	"github.com/millken/golog/internal/log"
	"github.com/millken/golog/internal/meta"
	"github.com/stretchr/testify/require"
)

func TestText(t *testing.T) {

	module := "test-module"
	level := log.INFO
	meta.ShowCallerInfo(module, level)
	text := NewText()
	_, err := text.Encode(nil)
	require.Error(t, err)
	e := &log.Entry{
		Timestamp: time.Now(),
		Module:    module,
		Level:     log.INFO,
		Message:   "test",
	}
	b, err := text.Encode(e)
	require.NoError(t, err)
	require.Contains(t, string(b), "test")
	require.Contains(t, string(b), "INF")
	require.Contains(t, string(b), "text_test.go")

}

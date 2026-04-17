package golog_test

import (
	"testing"

	"github.com/millken/golog"
	"github.com/stretchr/testify/require"
)

func TestTextEncoder(t *testing.T) {
	defer golog.DefaultCallerSkip.Store(3)

	module := "test-module"
	cs := golog.NewTextEncoder(golog.TextEncoderConfig{
		DisableTimestamp: true,
		DisableColor:     true,
	})
	golog.DefaultCallerSkip.Store(1)
	_, err := cs.Encode(nil)
	require.Error(t, err)
	e := &golog.Entry{
		Module:  module,
		Level:   golog.INFO,
		Message: "test",
	}
	e.SetFlag(golog.FlagCaller)

	b, err := cs.Encode(e)
	require.NoError(t, err)
	require.Contains(t, string(b), "test")
	require.Contains(t, string(b), "INF")
	require.Contains(t, string(b), "encoder_test.go")
}

func TestJSONEncoder(t *testing.T) {
	defer golog.DefaultCallerSkip.Store(3)

	module := "test-module"
	cs := golog.NewJSONEncoder(golog.JSONEncoderConfig{})
	_, err := cs.Encode(nil)
	require.Error(t, err)
	e := &golog.Entry{
		Module:  module,
		Level:   golog.INFO,
		Message: "test",
	}
	e.SetFlag(golog.FlagCaller)
	golog.DefaultCallerSkip.Store(1)
	b, err := cs.Encode(e)
	require.NoError(t, err)
	require.Contains(t, string(b), "test")
	require.Contains(t, string(b), "info")
	require.Contains(t, string(b), "encoder_test.go")
}

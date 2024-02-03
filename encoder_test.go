package golog_test

import (
	"testing"

	"github.com/millken/golog"
	"github.com/millken/x/buffer"
	"github.com/stretchr/testify/require"
)

func TestTextEncoder(t *testing.T) {
	require := require.New(t)
	cs := golog.NewEncoderText(golog.EncoderTextConfig{
		DisableTimestamp: true,
	})
	buf := buffer.Get()
	rec := golog.Record{
		Level:   golog.INFO,
		Message: "test",
		Config: golog.Config{
			Shortfile: true,
			Stack:     true,
		},
	}
	err := cs.Encode(buf, rec)
	require.NoError(err)

	require.Contains(string(buf.Bytes()), "test")
	require.Contains(string(buf.Bytes()), "info")
	require.Contains(string(buf.Bytes()), "encoder_test.go")
}

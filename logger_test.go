package golog_test

import (
	"testing"

	"github.com/millken/golog"
	"github.com/millken/golog/handlers/simple"
	"github.com/millken/x/buffer"
)

func TestLogger(t *testing.T) {
	// require := require.New(t)
	var buf buffer.Buffer
	l := golog.New()
	h1 := simple.NewSimple(&buf)
	golog.SetEncoder(golog.NewEncoderJSON(golog.EncoderJSONConfig{}))
	l.AddHandler(h1)
	l.Debug("debugmessage3", "a", 1, "b", "c")
	print(buf.String())
}

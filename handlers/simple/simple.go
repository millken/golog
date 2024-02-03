package simple

import (
	"io"

	"github.com/millken/golog"
)

type Simple struct {
	golog.HandlerBase
	io.Writer
}

func NewSimple(writer io.Writer) *Simple {
	return &Simple{Writer: writer}
}

func (h *Simple) SetWriter(w io.Writer) {
	h.Writer = w
}

func (h *Simple) Write(p []byte) (n int, err error) {
	return h.Writer.Write(p)
}

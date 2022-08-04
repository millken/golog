package handler

import (
	"io"

	"github.com/millken/golog/internal/log"
)

var (
	_ log.Handler = (*Writer)(nil)
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Handle(e *log.Entry) error {
	_, err := w.writer.Write(e.Bytes())
	return err
}

package writer

import (
	"io"

	"github.com/millken/golog"
	"github.com/millken/golog/handler"
)

type Handler struct {
	handler.AbstractHandler
	Output io.Writer
}

func (h *Handler) SetOutput(output io.Writer) {
	h.Output = output
}

func (h *Handler) Handle(entry *golog.Entry) error {
	_, err := h.Output.Write(entry.Data)
	return err
}

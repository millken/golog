package golog

import "io"

type FileHandler struct {
	AbstractHandler
	Output io.Writer
}

func (h *FileHandler) SetOutput(output io.Writer) {
	h.Output = output
}

func (h *FileHandler) Handle(entry *Entry) error {
	_, err := h.Output.Write(entry.Formatted)
	return err
}

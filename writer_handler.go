package golog

import "io"

// WriterHandler is a handler that writes log entries to an io.Writer.
type WriterHandler struct {
	AbstractHandler
	Output io.Writer
}

// SetOutput sets the output of the handler.
func (h *WriterHandler) SetOutput(output io.Writer) {
	h.Output = output
}

// Handle writes the log entry to the output.
func (h *WriterHandler) Handle(entry *Entry) error {
	_, err := h.Output.Write(entry.Data)
	return err
}

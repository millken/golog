package golog

import "io"

// Handler is the interface that must be implemented
type Handler interface {
	Handle(*Entry) error
	SetLevel(Level)
	Level() Level
	SetFormatter(Formatter)
	Formatter() Formatter
	SetDisableLogFields(bool)
	DisableLogFields() bool
	SetLevels(...Level)
	Levels() Levels
}

// baseHandler is the base type for all handlers.
type baseHandler struct {
	disableFields bool
	level         Level
	formatter     Formatter
	levels        Levels
}

// SetLevel sets the level of the handler.
func (h *baseHandler) SetLevel(level Level) {
	h.level = level
}

// Level returns the level of the handler.
func (h *baseHandler) Level() Level {
	return h.level
}

// SetFormatter sets the formatter of the handler.
func (h *baseHandler) SetFormatter(formatter Formatter) {
	h.formatter = formatter
}

// Formatter is the interface that must be implemented
func (h *baseHandler) Formatter() Formatter {
	return h.formatter
}

// SetDisableLogFields sets whether the handler is disabled to log fields.
func (h *baseHandler) SetDisableLogFields(disable bool) {
	h.disableFields = disable
}

// DisableLogFields returns whether the handler is disabled to log fields.
func (h *baseHandler) DisableLogFields() bool {
	return h.disableFields
}

// SetLevels sets the levels of the handler.
func (h *baseHandler) SetLevels(levels ...Level) {
	h.levels = levels
}

// Levels is a set of levels.
func (h *baseHandler) Levels() Levels {
	return h.levels
}

type loggerHandler struct {
	baseHandler
	writer io.Writer
}

// Handle writes the log entry to the output.
func (h *loggerHandler) Handle(entry *Entry) error {
	_, err := h.writer.Write(entry.Data)
	return err
}

// NewLoggerHandler creates a new logger handler.
func NewLoggerHandler(writer io.Writer) Handler {
	return &loggerHandler{
		writer: writer,
	}
}

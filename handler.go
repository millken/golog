package golog

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

// AbstractHandler is the base type for all handlers.
type AbstractHandler struct {
	disableFields bool
	level         Level
	formatter     Formatter
	levels        Levels
}

// SetLevel sets the level of the handler.
func (h *AbstractHandler) SetLevel(level Level) {
	h.level = level
}

// Level returns the level of the handler.
func (h *AbstractHandler) Level() Level {
	return h.level
}

// SetFormatter sets the formatter of the handler.
func (h *AbstractHandler) SetFormatter(formatter Formatter) {
	h.formatter = formatter
}

// Formatter is the interface that must be implemented
func (h *AbstractHandler) Formatter() Formatter {
	return h.formatter
}

// SetDisableLogFields sets whether the handler is disabled to log fields.
func (h *AbstractHandler) SetDisableLogFields(disable bool) {
	h.disableFields = disable
}

// DisableLogFields returns whether the handler is disabled to log fields.
func (h *AbstractHandler) DisableLogFields() bool {
	return h.disableFields
}

// SetLevels sets the levels of the handler.
func (h *AbstractHandler) SetLevels(levels ...Level) {
	h.levels = levels
}

// Levels is a set of levels.
func (h *AbstractHandler) Levels() Levels {
	return h.levels
}

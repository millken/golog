package golog

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

type AbstractHandler struct {
	disableFields bool
	level         Level
	formatter     Formatter
	levels        Levels
}

func (h *AbstractHandler) SetLevel(level Level) {
	h.level = level
}

func (h *AbstractHandler) Level() Level {
	return h.level
}
func (h *AbstractHandler) SetFormatter(formatter Formatter) {
	h.formatter = formatter
}

func (h *AbstractHandler) Formatter() Formatter {
	return h.formatter
}

func (h *AbstractHandler) SetDisableLogFields(disable bool) {
	h.disableFields = disable
}

func (h *AbstractHandler) DisableLogFields() bool {
	return h.disableFields
}

func (h *AbstractHandler) SetLevels(levels ...Level) {
	h.levels = levels
}

func (h *AbstractHandler) Levels() Levels {
	return h.levels
}

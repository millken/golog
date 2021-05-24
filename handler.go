package golog

type Handler interface {
	Handle(*Entry) error
	SetLevel(Level)
	GetLevel() Level
	SetFormatter(Formatter)
	GetFormatter() Formatter
}
type AbstractHandler struct {
	level     Level
	formatter Formatter
}

func (h *AbstractHandler) SetLevel(level Level) {
	h.level = level
}

func (h *AbstractHandler) GetLevel() Level {
	return h.level
}
func (h *AbstractHandler) SetFormatter(formatter Formatter) {
	h.formatter = formatter
}

func (h *AbstractHandler) GetFormatter() Formatter {
	return h.formatter
}

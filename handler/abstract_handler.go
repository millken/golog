package handler

import "github.com/millken/golog"

type AbstractHandler struct {
	disableFields bool
	level         golog.Level
	formatter     golog.Formatter
	levels        golog.Levels
}

func (h *AbstractHandler) SetLevel(level golog.Level) {
	h.level = level
}

func (h *AbstractHandler) Level() golog.Level {
	return h.level
}
func (h *AbstractHandler) SetFormatter(formatter golog.Formatter) {
	h.formatter = formatter
}

func (h *AbstractHandler) Formatter() golog.Formatter {
	return h.formatter
}

func (h *AbstractHandler) SetDisableLogFields(disable bool) {
	h.disableFields = disable
}

func (h *AbstractHandler) DisableLogFields() bool {
	return h.disableFields
}
func (h *AbstractHandler) SetLevels(levels ...golog.Level) {
	h.levels = levels
}

func (h *AbstractHandler) Levels() golog.Levels {
	return h.levels
}

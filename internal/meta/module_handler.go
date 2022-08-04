package meta

import (
	"os"

	"github.com/millken/golog/internal/handler"
	"github.com/millken/golog/internal/log"
)

var (
	defaultHandler = handler.NewWriter(os.Stdout)
)

func newModuledHandlers() *moduleHandlers {
	return &moduleHandlers{handlers: make(map[string]log.Handler)}
}

// moduleHandlers maintains log handlers based on modules.
type moduleHandlers struct {
	handlers map[string]log.Handler
}

// GetHandler returns the log handler for given module.
func (h *moduleHandlers) GetHandler(module string) log.Handler {
	handler, exists := h.handlers[module]
	if !exists {
		handler, exists = h.handlers[defaultModuleName]
		if !exists {
			return defaultHandler
		}
	}

	return handler
}

// SetDefaultHandler sets the default log handler.
func (h *moduleHandlers) SetDefaultHandler(handler log.Handler) {
	h.handlers[defaultModuleName] = handler
}

// SetHandler sets the log handler for given module.
func (h *moduleHandlers) SetHandler(module string, handler log.Handler) {
	h.handlers[module] = handler
}

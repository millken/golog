package golog

import (
	"io"
	"sync/atomic"
)

type Handler interface {
	Config() Config
	Encoder() Encoder
	SetWriter(io.Writer)
	io.Writer
}
type HandlerBase struct {
}

func (h *HandlerBase) Config() Config {
	return defaultConfigValue.Load().(Config)
}

func (h *HandlerBase) Encoder() Encoder {
	return defaultEncoder()
}

var (
	_                   Handler = (*baseHandler)(nil)
	defaultHandlerValue atomic.Value
)

func defaultHandler() Handler {
	return defaultHandlerValue.Load().(*baseHandler)
}

type baseHandler struct {
	HandlerBase
	Writer io.Writer
}

func newDefaultHandler(w io.Writer) *baseHandler {
	return &baseHandler{Writer: w}
}

func (h *baseHandler) SetWriter(w io.Writer) {
	h.Writer = w
}

func (h *baseHandler) Write(p []byte) (n int, err error) {
	return h.Writer.Write(p)
}

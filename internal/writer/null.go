package writer

import "github.com/millken/golog/internal/log"

var (
	_ log.Handler = (*Null)(nil)
)

type Null struct {
}

func NewNull() *Null {
	return &Null{}
}

func (w *Null) Handle(e *log.Entry) error {
	return nil
}

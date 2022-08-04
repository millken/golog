package writer

import (
	"io"
	"os"

	"github.com/millken/golog/internal/config"
)

var (
	_ io.Writer = (*File)(nil)
)

type File struct {
	cfg    config.FileConfig
	writer io.Writer
}

func NewFile(cfg config.FileConfig) (*File, error) {
	var writer io.Writer
	switch cfg.Path {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	case "":
		writer = io.Discard
	default:
		f, err := os.OpenFile(cfg.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		writer = f
	}

	return &File{
		cfg:    cfg,
		writer: writer,
	}, nil
}

func (w *File) Write(b []byte) (n int, err error) {
	return w.writer.Write(b)
}

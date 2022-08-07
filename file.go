package golog

import (
	"io"
	"os"
)

var (
	_ io.Writer = (*File)(nil)
)

// File is an implementation of io.Writer interface.
type File struct {
	cfg    FileConfig
	writer io.Writer
}

// NewFile creates and returns a new File writer.
func NewFile(cfg FileConfig) (*File, error) {
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

// Write writes the contents of b to the file.
func (w *File) Write(b []byte) (n int, err error) {
	return w.writer.Write(b)
}

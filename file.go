package golog

import (
	"bufio"
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
		f, err := os.OpenFile(cfg.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		writer = bufio.NewWriterSize(f, 4096)
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

// Flush flushes any buffered data to the underlying writer.
func (w *File) Flush() error {
	if f, ok := w.writer.(*bufio.Writer); ok {
		return f.Flush()
	}
	return nil
}

// Close flushes any buffered data and closes the underlying writer if it implements io.Closer.
func (w *File) Close() error {
	_ = w.Flush()
	if c, ok := w.writer.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

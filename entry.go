package golog

import (
	"io"
	"runtime"
	"sync"
	"time"
)

type Entry struct {
	Data      []byte
	Formatted []byte
	Timestamp time.Time
	Level     Level
	Fields    fields
}

var (
	entryPool = &sync.Pool{
		New: func() interface{} {
			return new(Entry)
		},
	}
)

func acquireEntry() *Entry {
	return entryPool.Get().(*Entry)
}

func releaseEntry(e *Entry) {
	e.Data = e.Data[:0]
	e.Formatted = e.Formatted[:0]
	e.Fields = e.Fields[:0]
	entryPool.Put(e)
}

func (e *Entry) Bytes() []byte {
	return e.Formatted
}

func (e *Entry) WriteByte(c byte) error {
	e.Formatted = append(e.Formatted, c)
	return nil
}

func (e *Entry) Write(p []byte) (int, error) {
	e.Formatted = append(e.Formatted, p...)
	return len(p), nil
}

func (e *Entry) WriteString(s string) (int, error) {
	e.Formatted = append(e.Formatted, s...)
	return len(s), nil
}

func (e *Entry) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(e.Formatted)
	return int64(n), err
}

func (e *Entry) Reset() {
	e.Formatted = e.Formatted[:0]
}

func (e *Entry) GetCaller(skip int) (string, int) {
	_, file, line, _ := runtime.Caller(skip)
	return file, line
}

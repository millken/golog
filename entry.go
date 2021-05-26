package golog

import (
	"io"
	"runtime"
	"sync"
	"time"
)

type Entry struct {
	Message   string
	Data      []byte
	Timestamp time.Time
	Level     Level
	Fields    []field
}

var (
	entryPool = &sync.Pool{
		New: func() interface{} {
			return &Entry{
				Data: make([]byte, 0, 500),
			}
		},
	}
)

func init() {
	entry := acquireEntry()
	releaseEntry(entry)
}

func acquireEntry() *Entry {
	return entryPool.Get().(*Entry)
}

func releaseEntry(e *Entry) {
	e.Message = e.Message[:0]
	e.Data = e.Data[:0]
	e.Fields = e.Fields[:0]
	entryPool.Put(e)
}

func (e *Entry) Bytes() []byte {
	return e.Data
}

func (e *Entry) WriteByte(c byte) error {
	e.Data = append(e.Data, c)
	return nil
}

func (e *Entry) Write(p []byte) (int, error) {
	e.Data = append(e.Data, p...)
	return len(p), nil
}

func (e *Entry) WriteString(s string) (int, error) {
	e.Data = append(e.Data, s...)
	return len(s), nil
}

func (e *Entry) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(e.Data)
	return int64(n), err
}

func (e *Entry) Reset() {
	e.Data = e.Data[:0]
}

func (e *Entry) GetCaller(skip int) (string, int) {
	_, file, line, _ := runtime.Caller(skip)
	return file, line
}

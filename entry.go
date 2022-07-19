package golog

import (
	"io"
	"sync"
	"time"
)

// Entry represents a log entry.
type Entry struct {
	Message    string
	Data       []byte
	Timestamp  time.Time
	Level      Level
	Fields     []Field
	fieldsLen  int
	callerSkip int
	caller     string
}

var (
	entryPool = &sync.Pool{
		New: func() interface{} {
			return &Entry{
				Data:   make([]byte, 0, 4096),
				Fields: make([]Field, 0, 512),
			}
		},
	}
)

func acquireEntry() *Entry {
	return entryPool.Get().(*Entry) //nolint:errcheck
}

func releaseEntry(e *Entry) {
	e.Message = e.Message[:0]
	e.Data = e.Data[:0]
	e.Fields = e.Fields[:0]
	e.fieldsLen = 0
	e.callerSkip = 0
	e.caller = e.caller[:0]
	entryPool.Put(e)
}

// Bytes returns the entry data as bytes.
func (e *Entry) Bytes() []byte {
	return e.Data
}

// WriteByte appends the byte to the entry data.
func (e *Entry) WriteByte(c byte) error {
	e.Data = append(e.Data, c)

	return nil
}

// Write appends the contents of p to the entry data.
func (e *Entry) Write(p []byte) (int, error) {
	e.Data = append(e.Data, p...)

	return len(p), nil
}

// WriteString appends the string to the entry data.
func (e *Entry) WriteString(s string) (int, error) {
	e.Data = append(e.Data, s...)
	return len(s), nil
}

// WriteTo writes the entry data to w.
func (e *Entry) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(e.Data)
	return int64(n), err
}

// FieldsLength returns the number of fields.
func (e *Entry) FieldsLength() int {
	return e.fieldsLen
}

// Reset resets the entry data.
func (e *Entry) Reset() {
	e.Data = e.Data[:0]
}

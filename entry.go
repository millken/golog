package golog

import (
	"io"

	"github.com/millken/gosync"
)

type Flag uint8

const (
	FlagNoColor Flag = 1 << iota
	FlagTime
	FlagCaller
	FlagStacktrace
	FlagName
	_
	_
	_
)

type Entry struct {
	Module     string
	Message    string
	Data       []byte
	Level      Level
	Fields     []Field
	fieldsLen  int
	callerSkip int
	caller     string
	flag       Flag
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

// SetFieldsLen sets the number of fields.
func (e *Entry) SetFieldsLen(n int) {
	e.fieldsLen = n
}

// CallerSkip returns the caller skip.
func (e *Entry) CallerSkip() int {
	return e.callerSkip
}

// SetCallerSkip sets the caller skip.
func (e *Entry) SetCallerSkip(n int) {
	e.callerSkip = n
}

// SetCaller sets the caller.
func (e *Entry) SetCaller(caller string) {
	e.caller = caller
}

// GetCaller gets the caller.
func (e *Entry) GetCaller() string {
	return e.caller
}

// SetFlag sets the flag.
func (e *Entry) SetFlag(flag Flag) {
	e.flag |= flag
}

// HasFlag returns true if the flag is set.
func (e *Entry) HasFlag(flag Flag) bool {
	return e.flag&flag == flag
}

// Reset resets the entry data.
func (e *Entry) Reset() {
	e.Data = e.Data[:0]
}

var (
	entryPool = gosync.NewPool(
		func() *Entry {
			return &Entry{
				Data:   make([]byte, 0, 4096),
				Fields: make([]Field, 0, 500),
			}
		})
)

// acquireEntry returns a new entry.
func acquireEntry() *Entry {
	return entryPool.Get() //nolint:errcheck
}

// releaseEntry releases the entry.
func releaseEntry(e *Entry) {
	e.Message = e.Message[:0]
	e.Data = e.Data[:0]
	e.Fields = e.Fields[:0]
	e.fieldsLen = 0
	e.callerSkip = 0
	e.caller = e.caller[:0]
	e.flag = 0
	entryPool.Put(e)
}

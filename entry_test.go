package golog

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEntry(t *testing.T) {
	require := require.New(t)
	entry := &Entry{
		Level: DEBUG,
	}
	require.Equal(DEBUG, entry.Level)
	require.Equal("", entry.Message)
	require.Equal([]byte(nil), entry.Data)
	require.Equal([]Field(nil), entry.Fields)
	require.Equal(0, entry.FieldsLength())
	n, err := entry.Write([]byte("hello"))
	require.Equal(5, n)
	require.Nil(err)
	require.Equal([]byte("hello"), entry.Data)

	entry.Message = "hello"
	require.Equal("hello", entry.Message)
	require.NoError(entry.WriteByte(' '))
	require.Equal(byte(' '), entry.Data[5])

	_, _ = entry.Write([]byte("world"))
	require.Equal([]byte("hello world"), entry.Data)

	_, _ = entry.WriteString("!")
	require.Equal([]byte("hello world!"), entry.Data)

	var a bytes.Buffer
	_, _ = entry.WriteTo(&a)
	require.Equal([]byte("hello world!"), a.Bytes())
	require.Equal(entry.Bytes(), a.Bytes())
}

package golog

import (
	"errors"
	"io"
	"testing"
)

func TestStdLog(t *testing.T) {
	Debug("std debug message")
	Info("std info message")
	Warn("std warning message")
	Error("std error message")
	StdSetLevel(DebugLevel)
	Debug("std debug message", Field("err", errors.New("error")))
	Info("std debug message with 2 fields", Field("a", 1), Field("b", true))
	StdSetOutput(io.Discard)
	Debug("std debug message")
}

func BenchmarkStdlog(b *testing.B) {
	StdSetOutput(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("abcde1234")
	}
}

func BenchmarkStdlogWithFields(b *testing.B) {
	StdSetOutput(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("abcde1234", Field("a", 1), Field("b", true))
	}
}

package golog

import (
	"io"
	"testing"
)

func TestStdLog(t *testing.T) {
	Debug("std debug message")
	Info("std info message")
	Warn("std warning message")
	StdEnableCaller = true
	Error("std error message")
	SetLevel(DebugLevel)
	StdEnableCaller = false
	Debug("std debug message")
	Info("std debug message with 2 fields", Field("a", 1), Field("b", true))
	SetOutput(io.Discard)
	Debug("std debug message")
}

func BenchmarkStdlog(b *testing.B) {
	SetOutput(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("abcde1234")
	}
}

func BenchmarkStdlogWithFields(b *testing.B) {
	SetOutput(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("abcde1234", Field("a", 1), Field("b", true))
	}
}

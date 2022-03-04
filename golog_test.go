package golog

import (
	"errors"
	"io"
	"testing"
)

func TestStdLog(t *testing.T) {
	Debugf("std debug message")
	Infof("std info message")
	Warnf("std warning message")
	Errorf("std error message")
	StdSetLevel(DebugLevel)
	WithField("err", errors.New("error")).Debugf("std debug message")
	WithFields(Field("a", 1), Field("b", true)).Infof("std debug message with %d fields", 2)
	StdSetOutput(io.Discard)
	Debugf("std debug message")
}

func BenchmarkStdlog(b *testing.B) {
	StdSetOutput(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Infof("abcde1234")
	}
}

func BenchmarkStdlogWithFields(b *testing.B) {
	StdSetOutput(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithFields(Field("a", 1), Field("b", true)).Infof("abcde1234")
	}
}

package golog

import (
	"errors"
	"io"
	"testing"
)

func TestStdLog(t *testing.T) {
	stdLog := NewStdLog()
	stdLog.Debugf("std debug message")
	stdLog.Infof("std info message")
	stdLog.Warnf("std warning message")
	stdLog.Errorf("std error message")
	stdLog.SetLevel(DebugLevel)
	stdLog.WithField("err", errors.New("error")).Debugf("std debug message")
	stdLog.WithFields(Field("a", 1), Field("b", true)).Infof("std info message with %d fields", 2)
	stdLog.SetOutput(io.Discard)
	stdLog.Debugf("std debug message")
}

func BenchmarkStdlog(b *testing.B) {
	stdLog := NewStdLog()
	stdLog.SetOutput(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdLog.Infof("abcde1234")
	}
}

func BenchmarkStdlogWithFields(b *testing.B) {
	stdLog := NewStdLog()
	stdLog.SetOutput(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdLog.WithFields(Field("a", 1), Field("b", true)).Infof("abcde1234")
	}
}

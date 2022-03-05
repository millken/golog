package golog

import (
	"errors"
	"io"
	"testing"
)

func TestGlobalLog(t *testing.T) {
	Debugf("std debug message")
	Infof("std info message")
	Warnf("std warning message")
	Errorf("std error message")
	WithField("err", errors.New("error")).Debugf("std debug message")
	WithFields(Field("a", 1), Field("b", true)).Infof("std info message with %d fields", 2)
	SetLevel(DebugLevel)
	Debugf("std debug message")
}

func BenchmarkGlobalLogger(b *testing.B) {
	SetOutput(io.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Infof("abcde1234")
	}
}

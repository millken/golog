package golog

import (
	"errors"
	"io"
	"os"
	"testing"
)

func TestGlobalLog(t *testing.T) {
	Debugf("std debug message")
	Infof("std info message")
	opt := StdOption{
		Output:           os.Stderr,
		DisableTimestamp: true,
		NoColor:          true,
	}
	stdLog := NewStdLog(opt)
	ReplaceGlobals(stdLog)
	Warnf("std warning message")
	Errorf("std error message")
	WithField("err", errors.New("error")).Debugf("std debug message")
	WithFields(Field("a", 1), Field("b", true)).Infof("std info message with %d fields", 2)
	Debugf("std debug message")
}

func BenchmarkGlobalLogger(b *testing.B) {
	opt := StdOption{
		Output:           io.Discard,
		DisableTimestamp: true,
		NoColor:          true,
	}
	stdLog := NewStdLog(opt)
	ReplaceGlobals(stdLog)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Infof("abcde1234")
	}
}

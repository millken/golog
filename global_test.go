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
		DisableTimestamp: false,
		NoColor:          true,
		EnableCaller:     true,
	}
	stdLog := NewStdLog(opt)
	ReplaceGlobals(stdLog)
	Warnf("std warning message")
	Errorf("std error message")
	WithField("err", errors.New("error")).Debugf("std debug message")
	WithFields(F("a", 1), F("b", true)).Infof("std info message with %d fields", 2)
	Debugf("std debug message")
}

func TestParseLevel(t *testing.T) {
	for _, test := range []struct {
		in   string
		want Level
	}{
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"error", ErrorLevel},
		{"fatal", FatalLevel},
		{"", Disabled},
		{"foo", Disabled},
	} {
		if got, err := ParseLevel(test.in); err == nil && got != test.want {
			t.Errorf("ParseLevel(%q) = %v, want %v", test.in, got, test.want)
		}
	}
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

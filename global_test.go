package golog

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestGlobal_Panic(t *testing.T) {
	require := require.New(t)
	var buf bytes.Buffer
	opt := StdOption{
		Output:           &buf,
		DisableTimestamp: true,
		NoColor:          true,
	}
	stdLog := NewStdLog(opt)

	ReplaceGlobals(stdLog)
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		Panicf("panic message")
	}()
	require.NotNil(recovered)
	require.Equal("PANIC panic message\n", buf.String())
	require.Equal("panic message", recovered)
}

func TestGlobal_Fatal(t *testing.T) {
	var buf bytes.Buffer
	opt := StdOption{
		Output:           &buf,
		DisableTimestamp: true,
		NoColor:          true,
	}
	stdLog := NewStdLog(opt)
	ReplaceGlobals(stdLog)

	if os.Getenv("BE_FATAL") == "1" {
		Fatalf("%s", "fatal")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestGlobal_Fatal")
	cmd.Env = append(os.Environ(), "BE_FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
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

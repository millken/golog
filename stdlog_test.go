package golog

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStdLog(t *testing.T) {
	opt := StdOption{
		EnableCaller: true,
		NoColor:      true,
	}
	stdLog := NewStdLog(opt)
	stdLog.Debug("std debug raw message")
	stdLog.Info("std info raw message")
	stdLog.Warn("std warning raw message")
	stdLog.Error("std error raw message")
	stdLog.Debugf("std debug message")
	stdLog.Infof("std info message")
	stdLog.Warnf("std warning message")
	stdLog.Errorf("std error message")
	stdLog.WithField("err", errors.New("error")).Debugf("std debug message")
	fields := []Field{
		F("int8", int8(1)),
		F("int16", int16(1)),
		F("int32", int32(1)),
		F("int64", int64(1)),
		F("uint8", uint8(1)),
		F("uint16", uint16(1)),
		F("uint32", uint32(1)),
		F("uint64", uint64(1)),
		F("float32", float32(1)),
		F("float64", float64(1)),
		F("bytes", []byte("bytes")),
		F("time", time.Now()),
		F("duration", time.Duration(time.Second*365000)),
		F("a", 1),
		F("b", true),
	}
	stdLog.WithFields(fields...).Infof("std info message with %d fields", len(fields))

	stdLog.Debugf("std debug message")
}

func TestStdLogRace(t *testing.T) {
	logger := NewStdLog()
	logger.Infof("should not race 01")
	go func() {
		logger.Infof("should not race 03")
	}()

	go func() {
		time.Sleep(200 * time.Microsecond)
		logger.Infof("should not race 04")
	}()
	time.Sleep(500 * time.Microsecond)
	logger.Infof("should not race 02")
}

func TestLogger_Panic(t *testing.T) {
	require := require.New(t)
	var buf bytes.Buffer
	opt := StdOption{
		Output:           &buf,
		DisableTimestamp: true,
		NoColor:          true,
	}
	stdLog := NewStdLog(opt)

	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		stdLog.Panicf("panic message")
	}()
	require.NotNil(recovered)
	require.Equal("PNC panic message\n", buf.String())
	require.Equal("panic message", recovered)
}

func TestLogger_Fatal(t *testing.T) {
	var buf bytes.Buffer
	opt := StdOption{
		Output:           &buf,
		DisableTimestamp: true,
		NoColor:          true,
	}
	stdLog := NewStdLog(opt)

	if os.Getenv("BE_FATAL") == "1" {
		stdLog.Fatalf("%s", "fatal")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLogger_Fatal")
	cmd.Env = append(os.Environ(), "BE_FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func BenchmarkStdlog(b *testing.B) {
	opt := StdOption{
		Output:           io.Discard,
		DisableTimestamp: true,
		NoColor:          true,
	}
	stdLog := NewStdLog(opt)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdLog.Infof("abcde1234")
	}
}

func BenchmarkStdlogWithFields(b *testing.B) {
	opt := StdOption{
		Output:           io.Discard,
		DisableTimestamp: true,
		NoColor:          true,
	}
	stdLog := NewStdLog(opt)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		stdLog.WithFields(F("a", 1), F("b", true)).Infof("abcde1234")
	}
}

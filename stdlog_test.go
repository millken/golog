package golog

import (
	"errors"
	"io"
	"testing"
	"time"
)

func TestStdLog(t *testing.T) {
	stdLog := NewStdLog()
	stdLog.Debugf("std debug message")
	stdLog.Infof("std info message")
	stdLog.Warnf("std warning message")
	stdLog.Errorf("std error message")
	stdLog.WithField("err", errors.New("error")).Debugf("std debug message")
	fields := []field{
		Field("int8", int8(1)),
		Field("int16", int16(1)),
		Field("int32", int32(1)),
		Field("int64", int64(1)),
		Field("uint8", uint8(1)),
		Field("uint16", uint16(1)),
		Field("uint32", uint32(1)),
		Field("uint64", uint64(1)),
		Field("float32", float32(1)),
		Field("float64", float64(1)),
		Field("bytes", []byte("bytes")),
		Field("time", time.Now()),
		Field("duration", time.Duration(time.Second*365000)),
		Field("a", 1),
		Field("b", true),
	}
	stdLog.WithFields(fields...).Infof("std info message with %d fields", 2)

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
		stdLog.WithFields(Field("a", 1), Field("b", true)).Infof("abcde1234")
	}
}

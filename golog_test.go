package golog

import (
	"syscall"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := NewLogger()
	logger.Debug("abdef", 1, true, "hell", time.Now())
}
func TestLoggerWithFields(t *testing.T) {
	logger := NewLogger()
	logger.WithFields(Fields{"a": 123, "b": true}).Debug("abdef")
}

/*
go test -benchmem -bench=. golog/*.go -memprofile profile_mem.out
go tool pprof golog.test profile_mem.out
*/
func BenchmarkLogger(b *testing.B) {
	logger := NewLogger()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("abcde1234")
	}
}

func BenchmarkLoggerWithFields(b *testing.B) {
	logger := NewLogger()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(Fields{"a": 123, "b": true}).Debug("abcde1234")

	}
}

func now() time.Time {
	var tv syscall.Timeval
	syscall.Gettimeofday(&tv)
	return time.Unix(0, syscall.TimevalToNsec(tv))
}

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

func BenchmarkNowGettimeofday(b *testing.B) {
	for i := 0; i < b.N; i++ {
		now()
	}
}

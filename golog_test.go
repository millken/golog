package golog

import (
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := NewLogger()
	logger.Debug("abdef", Field("a", 1), Field("b", true), Field("c", "hell"), Field("d", time.Now()))
}
func TestLoggerWithFields(t *testing.T) {
	logger := NewLogger()
	logger.Debug("abdef")
}

// /*
// go test -benchmem -bench=. golog/*.go -memprofile profile_mem.out
// go tool pprof golog.test profile_mem.out
// */
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
		logger.Debug("abcde1234", Field("a", 1), Field("b", true))

	}
}

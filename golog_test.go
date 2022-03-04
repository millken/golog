package golog

import (
	"os"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	logger := NewLogger()
	logger.Debug("abdef", Field("a", 1), Field("b", true), Field("c", "hell"), Field("d", time.Now()))
}
func TestLoggerWithFields(t *testing.T) {
	stdHandler = &FileHandler{
		Output: os.Stdout,
	}
	stdHandler.SetLevel(DebugLevel)
	stdFormatter = &TextFormatter{
		NoColor:              stdNoColor,
		TimeFormat:           stdTimeFormat,
		CallerSkipFrameCount: 6,
		EnableCaller:         true,
	}
	stdHandler.SetFormatter(stdFormatter)
	logger := NewLogger()
	logger.AddHandler(stdHandler)

	l := logger.WithFields(Field("a", 1), Field("b", true))
	l.Debugf("hello %s", "hell")
	l.Infof("hello %d", 435)
	l.WithField("c", "hell").Infof("hello 123")
	logger.Errorf("abcde1234")
}

// /*
// go test -benchmem -bench=. golog/*.go -memprofile profile_mem.out
// go tool pprof golog.test profile_mem.out
// */
func BenchmarkLoggerNoHandler(b *testing.B) {
	logger := NewLogger()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("abcde1234")
	}
}

func BenchmarkLoggerNoHandlerWithFields(b *testing.B) {
	logger := NewLogger()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("abcde1234", Field("a", 1), Field("b", true))

	}
}

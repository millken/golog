package golog

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logger := NewLogger()
	logger.WithFields(F("a", 1), F("b", true), F("c", "hell"), F("d", time.Now())).Debugf("abdef")
}

func TestLoggerWithOptions(t *testing.T) {
	require := require.New(t)
	logger := NewLogger()
	require.Equal(512, cap(logger.fields))
}

func TestLoggerWithFields(t *testing.T) {
	stdHandler := &WriterHandler{
		Output: os.Stdout,
	}
	stdHandler.SetLevel(DebugLevel)
	stdFormatter := &TextFormatter{
		NoColor:              false,
		TimeFormat:           StdTimeFormat,
		CallerSkipFrameCount: 6,
		EnableCaller:         true,
		PartsOrder:           []string{"time", "level", "caller", "message"},
	}
	stdHandler.SetFormatter(stdFormatter)
	stdHandler.SetDisableLogFields(true)
	logger := NewLogger()
	logger.AddHandler(stdHandler)

	l := logger.WithFields(F("a", 1), F("b", true))
	l.Debugf("hello %s", "hell")
	l.Infof("hello %d", 435)
	// l.WithField("c", "hell").Infof("hello 123")
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
		logger.Debugf("abcde1234")
	}
}

func BenchmarkLoggerNoHandlerWithFields(b *testing.B) {
	logger := NewLogger()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(F("a", 1), F("b", true)).Debugf("abcde1234")

	}
}

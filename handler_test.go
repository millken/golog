package golog

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriterHandler(t *testing.T) {
	fh := NewLoggerHandler(os.Stdout)
	fh.SetLevel(DebugLevel)
	formatter := NewTextFormatter()
	formatter.EnableCaller = true
	formatter.EnableStack = true
	fh.SetFormatter(formatter)
	logger := NewLogger()
	logger.AddHandler(fh)
	logger.Debugf("debug message")
	logger.Infof("info message")
	logger.Warnf("warning message")
	logger.Errorf("error message")
	logger.WithFields(F("a", 1), F("b", true)).Debugf("debug message with %d fields", 2)
	logger.WithField("c", "s").WithFields(F("a", 2), F("b", true)).Debugf("debug message with %d fields", 3)
}

func TestWriterHandlerWithJSONFormatter(t *testing.T) {
	fh := NewLoggerHandler(os.Stdout)
	fh.SetLevel(DebugLevel)
	fh.SetFormatter(&JSONFormatter{
		EnableCaller: true,
		EnableStack:  true,
	})
	logger := NewLogger()
	logger.AddHandler(fh)
	logger.Debugf("debug message")
	logger.Infof("info message")
	logger.Warnf("warning message")
	logger.Errorf("error message")
	logger.WithFields(F("a", 1), F("b", true)).Debugf("debug message")
	logger.WithField("c", "s").Debugf("debug message with %d fields", 2)

}

func TestDifferentLevelsGoToDifferentWriters(t *testing.T) {
	require := require.New(t)
	var a, b bytes.Buffer

	log := NewLogger()
	hand1 := NewLoggerHandler(&a)
	hand1.SetLevels(WarnLevel)
	hand1.SetFormatter(&TextFormatter{
		DisableTimestamp: true,
		PartsOrder:       []string{"level", "message"},
		NoColor:          true,
	})

	log.AddHandler(hand1)

	hand2 := NewLoggerHandler(&b)
	hand2.SetLevels(InfoLevel)
	hand2.SetFormatter(&TextFormatter{
		DisableTimestamp: true,
		PartsOrder:       []string{"level", "message"},
		NoColor:          true,
	})
	log.AddHandler(hand2)
	log.Warnf("send to a")
	log.Infof("send to b")

	require.Equal(a.String(), "WRN send to a\n")
	require.Equal(b.String(), "INF send to b\n")
}

func BenchmarkWriterHandler(b *testing.B) {
	fh := NewLoggerHandler(io.Discard)
	formatter := NewTextFormatter()
	formatter.EnableCaller = false
	formatter.DisableTimestamp = true
	fh.SetFormatter(formatter)
	logger := NewLogger()
	logger.AddHandler(fh)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debugf("abcde1234")
	}
}

func BenchmarkWriterHandlerWithFields(b *testing.B) {
	fh := NewLoggerHandler(io.Discard)

	formatter := NewTextFormatter()
	formatter.EnableCaller = false
	formatter.DisableTimestamp = true
	formatter.NoColor = true
	fh.SetFormatter(formatter)
	logger := NewLogger()
	logger.AddHandler(fh)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(F("a", 1), F("b", true)).Debugf("abcde1234")
	}
}

func BenchmarkJSONFormatterWriterHandler(b *testing.B) {
	fh := NewLoggerHandler(io.Discard)
	fh.SetFormatter(&JSONFormatter{
		DisableTimestamp: true,
		EnableCaller:     false,
	})
	logger := NewLogger()
	logger.AddHandler(fh)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debugf("abcde1234")
	}
}

func BenchmarkJSONFormatterWriterHandlerWithFields(b *testing.B) {
	fh := NewLoggerHandler(io.Discard)
	fh.SetFormatter(&JSONFormatter{
		DisableTimestamp: true,
		EnableCaller:     false,
	})
	logger := NewLogger()
	logger.AddHandler(fh)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(F("a", 1), F("b", true)).Debugf("abcde1234")

	}
}

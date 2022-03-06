package golog

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileHandler(t *testing.T) {
	fh := &FileHandler{
		Output: os.Stdout,
	}
	fh.SetLevel(DebugLevel)
	formatter := NewTextFormatter()
	formatter.EnableCaller = true
	fh.SetFormatter(formatter)
	logger := NewLogger()
	logger.AddHandler(fh)
	logger.Debugf("debug message")
	logger.Infof("info message")
	logger.Warnf("warning message")
	logger.Errorf("error message")
	logger.WithFields(Field("a", 1), Field("b", true)).Debugf("debug message with %d fields", 2)
}

func TestFileHandlerWithJSONFormatter(t *testing.T) {
	fh := &FileHandler{
		Output: os.Stdout,
	}
	fh.SetLevel(DebugLevel)
	fh.SetFormatter(&JSONFormatter{
		EnableCaller: true,
	})
	logger := NewLogger()
	logger.AddHandler(fh)
	logger.Debugf("debug message")
	logger.Infof("info message")
	logger.Warnf("warning message")
	logger.Errorf("error message")
	logger.WithFields(Field("a", 1), Field("b", true)).Debugf("debug message")

}

func TestDifferentLevelsGoToDifferentWriters(t *testing.T) {
	var a, b bytes.Buffer

	log := NewLogger()
	hand1 := &FileHandler{
		Output: &a,
	}
	hand1.SetLevels(WarnLevel)
	hand1.SetFormatter(&TextFormatter{
		DisableTimestamp: true,
		PartsOrder:       []string{"level", "message"},
		NoColor:          true,
	})

	log.AddHandler(hand1)

	hand2 := &FileHandler{
		Output: &b,
	}
	hand2.SetLevels(InfoLevel)
	hand2.SetFormatter(&TextFormatter{
		DisableTimestamp: true,
		PartsOrder:       []string{"level", "message"},
		NoColor:          true,
	})
	log.AddHandler(hand2)
	log.Warnf("send to a")
	log.Infof("send to b")

	assert.Equal(t, a.String(), "warn send to a\n")
	assert.Equal(t, b.String(), "info send to b\n")
}

func BenchmarkFileHandler(b *testing.B) {
	fh := &FileHandler{
		Output: io.Discard,
	}
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

func BenchmarkFileHandlerWithFields(b *testing.B) {
	fh := &FileHandler{
		Output: io.Discard,
	}

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
		logger.WithFields(Field("a", 1), Field("b", true)).Debugf("abcde1234")
	}
}

func BenchmarkJSONFormatterFileHandler(b *testing.B) {
	fh := &FileHandler{
		Output: io.Discard,
	}
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

func BenchmarkJSONFormatterFileHandlerWithFields(b *testing.B) {
	fh := &FileHandler{
		Output: io.Discard,
	}
	fh.SetFormatter(&JSONFormatter{
		DisableTimestamp: true,
		EnableCaller:     false,
	})
	logger := NewLogger()
	logger.AddHandler(fh)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(Field("a", 1), Field("b", true)).Debugf("abcde1234")

	}
}

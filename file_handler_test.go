package golog

import (
	"io"
	"os"
	"testing"
)

func TestFileHandler(t *testing.T) {
	fh := &FileHandler{
		Output: os.Stdout,
	}
	fh.SetLevel(DebugLevel)
	fh.SetFormatter(&TextFormatter{
		EnableCaller: true,
	})
	logger := NewLogger()
	logger.AddHandler(fh)
	logger.Debug(DebugLevel.String())
	logger.WithFields(Fields{"a": 123, "b": true}).Info(InfoLevel.String())
	logger.Warn(WarnLevel.String())
	logger.WithFields(Fields{"abcd": false}).Error(ErrorLevel.String())
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
	logger.Debug(DebugLevel.String())
	logger.WithFields(Fields{"a": 123, "b": true}).Info(InfoLevel.String())
	logger.Warn(WarnLevel.String())
	logger.WithFields(Fields{"abcd": false}).Error(ErrorLevel.String())
}

func BenchmarkFileHandler(b *testing.B) {
	fh := &FileHandler{
		Output: io.Discard,
	}
	fh.SetFormatter(&TextFormatter{})
	logger := NewLogger()
	logger.AddHandler(fh)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("abcde1234")
	}
}

func BenchmarkFileHandlerWithFields(b *testing.B) {
	fh := &FileHandler{
		Output: io.Discard,
	}
	fh.SetFormatter(&TextFormatter{})
	logger := NewLogger()
	logger.AddHandler(fh)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(Fields{"a": 123, "b": true}).Debug("abcde1234")

	}
}

func BenchmarkJSONFormatterFileHandler(b *testing.B) {
	fh := &FileHandler{
		Output: io.Discard,
	}
	fh.SetFormatter(&JSONFormatter{})
	logger := NewLogger()
	logger.AddHandler(fh)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("abcde1234")
	}
}

func BenchmarkJSONFormatterFileHandlerWithFields(b *testing.B) {
	fh := &FileHandler{
		Output: io.Discard,
	}
	fh.SetFormatter(&JSONFormatter{})
	logger := NewLogger()
	logger.AddHandler(fh)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.WithFields(Fields{"a": 123, "b": true}).Debug("abcde1234")

	}
}

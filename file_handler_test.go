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
		logger.Debugf("abcde1234")
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
		logger.WithFields(Field("a", 1), Field("b", true)).Debugf("abcde1234")

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
		logger.Debugf("abcde1234")
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
		logger.WithFields(Field("a", 1), Field("b", true)).Debugf("abcde1234")

	}
}

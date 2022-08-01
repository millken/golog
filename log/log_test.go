package log

import (
	"io"
	"os"
	"testing"

	"github.com/millken/golog"
)

func TestLog(t *testing.T) {
	f := golog.NewTextFormatter()
	f.EnableCaller = false
	f.NoColor = false
	f.EnableStack = false
	f.DisableTimestamp = true
	SetFormatter(f)
	SetOutput(os.Stdout)
	SetLevel(golog.DebugLevel)
	Debug("debug raw")
	Info("info raw")
	Warn("warn raw")
	Error("error raw")

	Debugf("debug")
	Infof("info")
	Warnf("warn")
	Errorf("error")
	WithField("a", 1).Debugf("debug")
	WithField("a", 1).WithField("b", true).Infof("debug multi field")
	WithFields(golog.F("a", 1), golog.F("b", true)).Debugf("debug with fields")
	l := WithField("a", 2)
	l.Infof("info")
	l1 := l.WithField("b", 3)
	l1.Warnf("warn")
	go func() {
		l1.Errorf("error")
		l1.WithFields(golog.F("a", 1), golog.F("b", true)).Errorf("error with field %d", 1)
	}()

	SetFormatter(&golog.JSONFormatter{})
	Debugf("debug")
	Infof("info")
	Warnf("warn")
	Errorf("error")
	WithField("a", 1).Debugf("debug")
	WithFields(golog.F("a", 1), golog.F("b", true)).Debugf("debug with fields")

}

func BenchmarkLogger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	SetOutput(io.Discard)
	f := golog.NewTextFormatter()
	f.EnableCaller = false
	f.NoColor = true
	f.EnableStack = false
	f.DisableTimestamp = true
	SetFormatter(f)
	for i := 0; i < b.N; i++ {
		Infof("abcde1234")
	}
}

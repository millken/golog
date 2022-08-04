package golog

import (
	"testing"

	"github.com/millken/golog/internal/config"
	"github.com/millken/golog/internal/log"
	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	require := require.New(t)
	require.NoError(LoadConfig("./internal/config/testdata/sample.yml"))
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
	WithField("a", 1).Debug("debug")
	WithField("a", 1).Info("info")
	WithField("a", 1).Warn("warn")
	WithField("a", 1).Error("error")
	WithFields(F("a", 1), F("b", 2)).Debug("debug")
	WithFields(F("a", 1), F("b", 2)).Info("info")
	WithFields(F("a", 1), F("b", 2)).Warn("warn")
	WithFields(F("a", 1), F("b", 2)).Error("error")
	l := WithFields(F("a", 1), F("b", 2))
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn")
	l.Error("error")
	ll := l.WithField("c", 3)
	ll.Debug("debug")
	ll.Info("info")
	ll.Warn("warn")
	ll.Error("error")
	l3 := New("test")
	l3.Infof("debug %s", "debug12")
	l3.Warnf("warn %s", "warn12")
	l3.Errorf("error %s", "error")
	l3.WithField("a", 1).Infof("info %s", "test")
	l3.WithField("a", 1).Warnf("warn %s", "test")
	l3.WithField("a", 1).Errorf("error %s", "test")
}

func TestDebugLog(t *testing.T) {
	require := require.New(t)
	require.NoError(LoadConfig("./internal/config/testdata/bench.yml"))
	Info("info")
}

func BenchmarkLog(b *testing.B) {
	require := require.New(b)
	cfg := config.Config{
		Level:    log.INFO,
		Encoding: "console",
		ConsoleEncoderConfig: config.ConsoleEncoderConfig{
			DisableTimestamp: true,
		},
		Writer: config.WriterConfig{
			Type: "file",
			FileConfig: config.FileConfig{
				Path: "",
			},
		},
	}
	log, err := NewLoggerByConfig("test", cfg)
	require.NoError(err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("info")
	}
}

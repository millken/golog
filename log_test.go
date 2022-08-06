package golog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	require.NoError(LoadConfig("./testdata/sample.yml"))
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
	WithField("a", 1).Debug("debug")
	WithField("a", 1).Info("info")
	WithField("a", 1).Warn("warn")
	WithField("a", 1).Error("error")

	f1 := Fields{"a": 1, "b": 2}
	WithFields(f1).Debug("debug")
	WithFields(f1).Info("info")
	WithFields(f1).Warn("warn")
	WithFields(f1).Error("error")
	l := WithFields(f1)
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
	l4 := l3.WithFields(f1)
	l4.Infof("info %s", "test")
	l4.WithField("c", 3).Warnf("warn %s", "test")
}

func TestDebugLog(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	require.NoError(LoadConfig("./testdata/bench.yml"))
	Info("info")
}

func BenchmarkLog_Console(b *testing.B) {
	require := require.New(b)
	cfg := Config{
		Level:    INFO,
		Encoding: "console",
		ConsoleEncoderConfig: ConsoleEncoderConfig{
			DisableTimestamp: true,
		},
		Writer: WriterConfig{
			Type: "file",
			FileConfig: FileConfig{
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

func BenchmarkLog_JSON(b *testing.B) {
	require := require.New(b)
	cfg := Config{
		Level:    INFO,
		Encoding: "json",
		JSONEncoderConfig: JSONEncoderConfig{
			DisableTimestamp: true,
		},
		Writer: WriterConfig{
			Type: "file",
			FileConfig: FileConfig{
				Path: "",
			},
		},
	}
	log, err := NewLoggerByConfig("test2", cfg)
	require.NoError(err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("info")
	}
}

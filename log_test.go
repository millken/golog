package golog

import (
	"testing"

	"errors"

	"github.com/stretchr/testify/require"
)

func makeFields() []interface{} {
	return []interface{}{
		"a", 1,
		"b", true,
		"c", 1.234,
		"d", "abc",
		"e", errors.New("abc"),
		"f", []byte("abc"),
		"g", []int{1, 2, 3},
		"h", []string{"a", "b", "c"},
		"i", []interface{}{1, "a", true},
		"j", map[string]interface{}{"a": 1, "b": "c"},
	}
}

func TestLog(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	require.NoError(LoadConfig("./testdata/sample.yml"))
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
	WithValues("a", 1).Debug("debug")
	WithValues("a", 1).Info("info")
	WithValues("a", 1).Warn("warn")
	Error("error", makeFields()...)

	f1 := []interface{}{"a", 1, "b", 2}
	WithValues(f1).Debug("debug")
	WithValues(f1).Info("info")
	WithValues(f1).Warn("warn")
	WithValues(f1).Error("error")
	l := WithValues(f1)
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn")
	l.Error("error")
	ll := l.WithValues("c", 3)
	ll.Debug("debug")
	ll.Info("info")
	ll.Warn("warn")
	ll.Error("error")
	l3 := New("test")
	l3.Infof("debug %s", "debug12")
	l3.Warnf("warn %s", "warn12")
	l3.Errorf("error %s", "error")
	l3.WithValues("a", 1).Infof("info %s", "test")
	l3.WithValues("a", 1).Warnf("warn %s", "test")
	l3.WithValues("a", 1).Errorf("error %s", "test")
	l4 := l3.WithValues(f1)
	l4.Infof("info %s", "test")
	l4.WithValues("c", 3).Warnf("warn %s", "test")
}

func TestLog_JSON(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	require.NoError(LoadConfig("./testdata/sample_json.yml"))
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error", makeFields()...)
	WithValues("a", 1).Debug("debug")
	WithValues("a", 1).Info("info")
	WithValues("a", 1).Warn("warn")
	WithValues("a", 1).Error("error")

	f1 := []interface{}{"a", 1, "b", 2}
	WithValues(f1).Debug("debug")
	WithValues(f1).Info("info")
	WithValues(f1).Warn("warn")
	WithValues(f1).Error("error")
	l := WithValues(f1)
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn")
	l.Error("error")
	ll := l.WithValues("c", 3)
	ll.Debug("debug")
	ll.Info("info")
	ll.Warn("warn")
	ll.Error("error", makeFields())
	l3 := New("test")
	l3.Infof("debug %s", "debug12")
	l3.Warnf("warn %s", "warn12")
	l3.Errorf("error %s", "error")
	l3.WithValues("a", 1).Infof("info %s", "test")
	l3.WithValues("a", 1).Warnf("warn %s", "test")
	l3.WithValues("a", 1).Errorf("error %s", "test")
	l4 := l3.WithValues(f1)
	l4.Infof("info %s", "test")
	l4.WithValues("c", 3).Warnf("warn %s", "test")
}

func TestDebugLog(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	require.NoError(LoadConfig("./testdata/debug.yml"))
	Info("The quick brown fox jumps over the lazy dog",
		"a", 1,
		"b", true,
		"c", 1.234,
	)
	Info("The quick brown fox jumps over the lazy dog",
		Fields{"a": 1,
			"b": true,
			"c": 1.234},
	)
}

func BenchmarkLogText(b *testing.B) {
	require := require.New(b)
	cfg := Config{
		Level:    INFO,
		Encoding: TextEncoding,
		TextEncoder: TextEncoderConfig{
			DisableTimestamp: true,
			DisableColor:     true,
		},
		Handler: HandlerConfig{
			Type: "file",
			File: FileConfig{
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

func BenchmarkLogText_WithField(b *testing.B) {
	require := require.New(b)
	cfg := Config{
		Level:    INFO,
		Encoding: TextEncoding,
		TextEncoder: TextEncoderConfig{
			DisableTimestamp: true,
			DisableColor:     true,
		},
		Handler: HandlerConfig{
			Type: "file",
			File: FileConfig{
				Path: "",
			},
		},
	}
	log, err := NewLoggerByConfig("test", cfg)
	require.NoError(err)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("The quick brown fox jumps over the lazy dog",
			"a", 1,
			"b", true,
			"c", 1.234,
		)
	}
}

func BenchmarkLogJSON(b *testing.B) {
	require := require.New(b)
	cfg := Config{
		Level:    INFO,
		Encoding: JSONEncoding,
		JSONEncoder: JSONEncoderConfig{
			DisableTimestamp: true,
		},
		Handler: HandlerConfig{
			Type: "file",
			File: FileConfig{
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

func BenchmarkLogJSON_WithField(b *testing.B) {
	require := require.New(b)
	cfg := Config{
		Level:    INFO,
		Encoding: JSONEncoding,
		JSONEncoder: JSONEncoderConfig{
			DisableTimestamp: true,
			ShowModuleName:   true,
		},
		Handler: HandlerConfig{
			Type: "file",
			File: FileConfig{
				Path: "",
			},
		},
	}
	log, err := NewLoggerByConfig("test", cfg)
	require.NoError(err)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info("The quick brown fox jumps over the lazy dog",
			"a", 1,
			"b", true,
			"c", 1.234,
		)
	}
}

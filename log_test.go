package golog_test

// import (
// 	"encoding/json"
// 	"testing"
// 	"time"

// 	"errors"

// 	"github.com/millken/golog"
// 	"github.com/millken/x/buffer"
// 	"github.com/stretchr/testify/require"
// )

// var (
// 	// fakeDisableColor is used to disable color output in tests.
// 	fakeDisableColor = true
// 	// fakeDisableTimestamp is used to disable timestamp output in tests.
// 	fakeDisableTimestamp = true
// 	// fakeDisableCaller is used to disable caller output in tests.
// 	fakeDisableCaller = true
// 	// fakeDisableStacktrace is used to disable stacktrace output in tests.
// 	fakeDisableStacktrace = true
// )

// func makeFields() []interface{} {
// 	return []interface{}{
// 		"a", 1,
// 		"b", true,
// 		"c", 1.234,
// 		"d", "abc",
// 		"e", errors.New("abc"),
// 		"f", []byte("abc"),
// 		"g", []int{1, 2, 3},
// 		"h", []string{"a", "b", "c"},
// 		"i", []interface{}{1, "a", true},
// 		"j", map[string]interface{}{"a": 1, "b": "c"},
// 	}
// }

// func validJSON(s string) bool {
// 	var js map[string]interface{}
// 	return json.Unmarshal([]byte(s), &js) == nil
// }

// func TestJSONLog(t *testing.T) {
// 	require := require.New(t)
// 	var buf buffer.Buffer
// 	// golog.SetWriter(&buf)
// 	// golog.SetLevel(golog.INFO)
// 	// golog.SetEncoding(golog.JSONEncoding)
// 	golog.Debug("debug")
// 	require.Empty(buf.String())
// 	golog.Info("info")
// 	require.True(validJSON(buf.String()))
// 	require.Contains(buf.String(), `"level":"info"`)
// 	require.Contains(buf.String(), `"message":"info"`)
// 	buf.Reset()
// 	golog.Warn("warn")
// 	golog.Error("error")
// 	buf.Reset()
// 	golog.WithValues("a", 1).Debug("debug")
// 	golog.WithValues("a", 1).Info("info")
// 	require.True(validJSON(buf.String()))
// 	require.Contains(buf.String(), `"a":1`)
// 	require.Contains(buf.String(), `"level":"info"`)
// 	require.Contains(buf.String(), `"message":"info"`)
// 	buf.Reset()
// 	golog.WithValues("a", 1).Warn("warn")
// 	golog.Error("error", makeFields()...)

// 	f1 := []interface{}{"a", 1, "b", 2}
// 	golog.WithValues(f1...).Debug("debug")
// 	golog.WithValues(f1...).Info("info")
// 	golog.WithValues(f1...).Warn("warn")
// 	golog.WithValues(f1...).Error("error")
// 	l := golog.WithValues(f1...)
// 	l.Debug("debug")
// 	l.Info("info")
// 	l.Warn("warn")
// 	l.Error("error")
// 	ll := l.WithValues("c", 3)
// 	ll.Debug("debug")
// 	ll.Info("info")
// 	ll.Warn("warn")
// 	buf.Reset()
// 	ll.Error("error")
// 	require.True(validJSON(buf.String()))
// 	require.Contains(buf.String(), "level\":\"error\",\"message\":\"error\",\"a\":1,\"b\":2,\"c\":3")
// 	l3 := golog.New("test")
// 	l3.Infof("debug %s", "debug12")
// 	l3.Warnf("warn %s", "warn12")
// 	l3.Errorf("error %s", "error")
// 	l3.WithValues("a", 1).Infof("info %s", "test")
// 	l3.WithValues("a", 1).Warnf("warn %s", "test")
// 	buf.Reset()
// 	l3.WithValues("a", 1).Errorf("error %s", "test")
// 	require.True(validJSON(buf.String()))
// 	require.Contains(buf.String(), "level\":\"error\",\"message\":\"error test\",\"a\":1")
// 	l4 := l3.WithValues(f1...)
// 	l4.Infof("info %s", "test")
// 	buf.Reset()
// 	l4.WithValues("c", 3).Warnf("warn %s", "test")
// 	require.True(validJSON(buf.String()))
// 	require.Contains(buf.String(), "level\":\"warning\",\"message\":\"warn test\",\"a\":1,\"b\":2,\"c\":3")
// }

// func TestLog_JSON_Output(t *testing.T) {
// 	// defer resetConfigs()
// 	// require := require.New(t)
// 	// require.NoError(golog.LoadConfig("./testdata/sample_json.yml"))
// 	golog.Debug("debug")
// 	golog.Info("info")
// 	golog.Warn("warn")
// 	golog.Info("test int", "int8", int8(1), "int16", int16(2), "int32", int32(3), "int64", int64(4))
// 	golog.Info("test uint", "uint", uint(0), "uint8", uint8(1), "uint16", uint16(2), "uint32", uint32(3), "uint64", uint64(4))
// 	golog.Info("test float", "float32", float32(1.1), "float64", float64(2.2))
// 	golog.Warn("test bool", "bool", true)
// 	golog.Warn("test string", "string", "string")
// 	golog.Warn("test error", "error", errors.New("error"))
// 	golog.Warn("test nil", "nil", nil)
// 	golog.Error("test map", "map", map[string]interface{}{"a": 1, "b": true})
// 	golog.Error("test array", "array", []interface{}{1, true, "string"})
// 	golog.Error("test struct", "struct", struct{}{})
// 	golog.Error("test time", "time", time.Now())
// 	golog.Error("test duration", "duration", time.Duration(1))
// 	golog.Error("test json.Number", "json.Number", json.Number("1.1"))
// 	golog.Error("error", makeFields()...)
// 	golog.WithValues("int8", int8(1), "int16", int16(2), "int32", int32(3), "int64", int64(4)).Debug("debug")
// 	golog.WithValues("uint", uint(0), "uint8", uint8(1), "uint16", uint16(2), "uint32", uint32(3), "uint64", uint64(4)).Info("info")
// 	golog.WithValues("float32", float32(1.1), "float64", float64(2.2)).Warn("warn")
// 	golog.WithValues("bool", true, "error", errors.New("error"), "nil", nil, "map", map[string]interface{}{"a": 1, "b": true}, "array", []interface{}{1, true, "string"}, "struct", struct{}{}, "time2", time.Now(), "duration", time.Duration(1)).Error("error")

// 	f1 := []interface{}{"a", 1, "b", 2}
// 	golog.WithValues(f1...).Debug("debug")
// 	golog.WithValues(f1...).Info("info")
// 	golog.WithValues(f1...).Warn("warn")
// 	golog.WithValues(f1...).Error("error")
// 	l := golog.WithValues(f1...)
// 	l.Debug("debug")
// 	l.Info("info")
// 	l.Warn("warn")
// 	l.Error("error")
// 	ll := l.WithValues("c", 3)
// 	ll.Debug("debug")
// 	ll.Info("info")
// 	ll.Warn("warn")
// 	ll.Error("error", makeFields())
// 	l3 := golog.New("test")
// 	l3.Infof("debug %s", "debug12")
// 	l3.Warnf("warn %s", "warn12")
// 	l3.Errorf("error %s", "error")
// 	l3.WithValues("a", 1).Infof("info %s", "test")
// 	l3.WithValues("a", 1).Warnf("warn %s", "test")
// 	l3.WithValues("a", 1).Errorf("error %s", "test")
// 	l4 := l3.WithValues(f1...)
// 	l4.Infof("info %s", "test")
// 	l4.WithValues("c", 3).Warnf("warn %s", "test")
// }

// func TestDebugLog(t *testing.T) {
// 	t.Skip()
// 	// require := require.New(t)
// 	// require.NoError(golog.LoadConfig("./testdata/debug.yml"))
// 	golog.Info("The quick brown fox jumps over the lazy dog",
// 		"a", 1,
// 		"b", true,
// 		"c", 1.234,
// 	)
// }

// func BenchmarkLogText(b *testing.B) {
// 	// require := require.New(b)
// 	// cfg := golog.Config{
// 	// 	Level: golog.INFO,
// 	// 	// Encoding: golog.TextEncoding,
// 	// 	// TextEncoder: golog.TextEncoderConfig{
// 	// 	// 	DisableTimestamp: fakeDisableTimestamp,
// 	// 	// 	DisableColor:     fakeDisableColor,
// 	// 	// },
// 	// 	// Handler: golog.HandlerConfig{
// 	// 	// 	Type: "file",
// 	// 	// 	File: golog.FileConfig{
// 	// 	// 		Path: "",
// 	// 	// 	},
// 	// 	// },
// 	// }
// 	// if !fakeDisableCaller {
// 	// 	cfg.CallerLevels = []golog.Level{golog.DEBUG, golog.INFO, golog.WARNING, golog.ERROR}
// 	// }
// 	// if !fakeDisableStacktrace {
// 	// 	cfg.StacktraceLevels = []golog.Level{golog.DEBUG, golog.INFO, golog.WARNING, golog.ERROR}
// 	// }
// 	// log, err := golog.NewLoggerByConfig("test1", cfg)
// 	// require.NoError(err)
// 	// b.ReportAllocs()
// 	// b.ResetTimer()
// }

// func BenchmarkLogText_WithField(b *testing.B) {
// 	// require := require.New(b)
// 	// cfg := golog.Config{
// 	// 	Level: golog.INFO,
// 	// 	// Encoding: golog.TextEncoding,
// 	// 	// TextEncoder: golog.TextEncoderConfig{
// 	// 	// 	DisableTimestamp: fakeDisableTimestamp,
// 	// 	// 	DisableColor:     fakeDisableColor,
// 	// 	// },
// 	// 	// Handler: golog.HandlerConfig{
// 	// 	// 	Type: "file",
// 	// 	// 	File: golog.FileConfig{
// 	// 	// 		Path: "",
// 	// 	// 	},
// 	// 	// },
// 	// }
// 	// if !fakeDisableCaller {
// 	// 	cfg.CallerLevels = []golog.Level{golog.DEBUG, golog.INFO, golog.WARNING, golog.ERROR}
// 	// }
// 	// if !fakeDisableStacktrace {
// 	// 	cfg.StacktraceLevels = []golog.Level{golog.DEBUG, golog.INFO, golog.WARNING, golog.ERROR}
// 	// }
// 	// log, err := golog.NewLoggerByConfig("test2", cfg)
// 	// require.NoError(err)

// 	// b.ReportAllocs()
// 	// b.ResetTimer()
// 	// for i := 0; i < b.N; i++ {
// 	// 	log.Info("The quick brown fox jumps over the lazy dog",
// 	// 		"a", 1,
// 	// 		"b", true,
// 	// 		"c", 1.234,
// 	// 	)
// 	// }
// }

// func BenchmarkLogJSON(b *testing.B) {
// 	// require := require.New(b)
// 	// cfg := golog.Config{
// 	// 	Level: golog.INFO,
// 	// 	// Encoding: golog.JSONEncoding,
// 	// 	// JSONEncoder: golog.JSONEncoderConfig{
// 	// 	// 	DisableTimestamp: fakeDisableTimestamp,
// 	// 	// },
// 	// 	// Handler: golog.HandlerConfig{
// 	// 	// 	Type: "file",
// 	// 	// 	File: golog.FileConfig{
// 	// 	// 		Path: "",
// 	// 	// 	},
// 	// 	// },
// 	// }
// 	// if !fakeDisableCaller {
// 	// 	cfg.CallerLevels = []golog.Level{golog.DEBUG, golog.INFO, golog.WARNING, golog.ERROR}
// 	// }
// 	// if !fakeDisableStacktrace {
// 	// 	cfg.StacktraceLevels = []golog.Level{golog.DEBUG, golog.INFO, golog.WARNING, golog.ERROR}
// 	// }
// 	// log, err := golog.NewLoggerByConfig("test3", cfg)
// 	// require.NoError(err)
// 	// b.ReportAllocs()
// 	// b.ResetTimer()
// 	// for i := 0; i < b.N; i++ {
// 	// 	log.Info("info")
// 	// }
// }

// func BenchmarkLogJSON_WithField(b *testing.B) {
// 	// require := require.New(b)
// 	// cfg := golog.Config{
// 	// 	Level: golog.INFO,
// 	// 	// Encoding: golog.JSONEncoding,
// 	// 	// JSONEncoder: golog.JSONEncoderConfig{
// 	// 	// 	DisableTimestamp: fakeDisableTimestamp,
// 	// 	// 	ShowModuleName:   false,
// 	// 	// },
// 	// 	// Handler: golog.HandlerConfig{
// 	// 	// 	Type: "file",
// 	// 	// 	File: golog.FileConfig{
// 	// 	// 		Path: "",
// 	// 	// 	},
// 	// 	// },
// 	// }
// 	// if !fakeDisableCaller {
// 	// 	cfg.CallerLevels = []golog.Level{golog.DEBUG, golog.INFO, golog.WARNING, golog.ERROR}
// 	// }
// 	// if !fakeDisableStacktrace {
// 	// 	cfg.StacktraceLevels = []golog.Level{golog.DEBUG, golog.INFO, golog.WARNING, golog.ERROR}
// 	// }
// 	// log, err := golog.NewLoggerByConfig("test4", cfg)
// 	// require.NoError(err)

// 	// b.ReportAllocs()
// 	// b.ResetTimer()
// 	// for i := 0; i < b.N; i++ {
// 	// 	log.Info("The quick brown fox jumps over the lazy dog",
// 	// 		"a", 1,
// 	// 		"b", true,
// 	// 		"c", 1.234,
// 	// 	)
// 	// }
// }

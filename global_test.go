package golog_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"os/exec"
	_ "runtime"
	"sync"
	"testing"
	"time"
	_ "unsafe"

	"github.com/millken/golog"
	"github.com/millken/golog/internal/buffer"
	"github.com/stretchr/testify/require"
)

func TestGlobalUsage(t *testing.T) {
	require := require.New(t)
	defer resetConfigs()
	var buf bytes.Buffer
	golog.SetLevel(golog.DEBUG)
	golog.SetCallerLevels(golog.DEBUG, golog.INFO, golog.WARNING, golog.ERROR, golog.FATAL, golog.PANIC)
	golog.SetStacktraceLevels(golog.PANIC, golog.FATAL, golog.ERROR, golog.WARNING)
	golog.SetEncoding(golog.TextEncoding)
	golog.SetTextEncoderConfig(golog.TextEncoderConfig{DisableTimestamp: true, DisableColor: true, ShowModuleName: true})
	golog.SetWriter(&buf)
	golog.Infof("hello %s", "world")
	require.Contains(buf.String(), "hello world")
	require.Contains(buf.String(), "global_test.go")
	buf.Reset()
	golog.Debug("test int", "int8", int8(1), "int16", int16(2), "int32", int32(3), "int64", int64(4))
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test int")
	require.Contains(buf.String(), "int8=1")
	require.Contains(buf.String(), "int16=2")
	require.Contains(buf.String(), "int32=3")
	require.Contains(buf.String(), "int64=4")
	buf.Reset()
	golog.Debug("test uint", "uint", uint(0), "uint8", uint8(1), "uint16", uint16(2), "uint32", uint32(3), "uint64", uint64(4))
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test uint")
	require.Contains(buf.String(), "uint=0")
	require.Contains(buf.String(), "uint8=1")
	require.Contains(buf.String(), "uint16=2")
	require.Contains(buf.String(), "uint32=3")
	require.Contains(buf.String(), "uint64=4")
	buf.Reset()
	golog.Debug("test float", "float32", float32(1.1), "float64", float64(2.2))
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test float")
	require.Contains(buf.String(), "float32=1.1")
	require.Contains(buf.String(), "float64=2.2")
	buf.Reset()
	golog.Debug("test bool", "bool", true)
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test bool")
	require.Contains(buf.String(), "bool=true")
	buf.Reset()
	golog.Debug("test string", "string", "string")
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test string")
	require.Contains(buf.String(), "string=string")
	buf.Reset()
	golog.Debug("test error", "error", errors.New("error"))
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test error")
	require.Contains(buf.String(), "error=error")
	buf.Reset()
	golog.Debug("test nil", "nil", nil)
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test nil")
	require.Contains(buf.String(), "nil=null")
	buf.Reset()
	golog.Debug("test map", "map", map[string]interface{}{"a": 1, "b": true})
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test map")
	require.Contains(buf.String(), "map={\"a\":1,\"b\":true}")
	buf.Reset()
	golog.Debug("test array", "array", []interface{}{1, true, "string"})
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test array")
	require.Contains(buf.String(), "array=[1,true,\"string\"]")
	buf.Reset()
	golog.Debug("test struct", "struct", struct{}{})
	require.Contains(buf.String(), "global_test.go")
	require.Contains(buf.String(), "test struct")
	require.Contains(buf.String(), "struct={}")
	buf.Reset()
	golog.Debug("test time", "time", time.Now())
	golog.Debug("test duration", "duration", time.Duration(1))
	golog.Warn("test json.Number", "json.Number", json.Number("1.1"))
	golog.Error("hello world with fields", "a", 1, "b", true, "c", "string")
	resetConfigs()
}

func TestGlobalLog(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	var buf buffer.Buffer
	golog.SetWriter(&buf)
	golog.Debug("debug message3")
	require.Empty(buf.String())
	golog.Info("info message3")
	require.Contains(buf.String(), "info message")
	buf.Reset()
	golog.Warn("warning message")
	require.Contains(buf.String(), "warning message")
	buf.Reset()
	golog.Error("error message")
	require.Contains(buf.String(), "error message")
	buf.Reset()
	golog.Debugf("debug message %s", "debug")
	require.Empty(buf.String())
	buf.Reset()
	golog.Infof("info message %s", "info")
	require.Contains(buf.String(), "info message info")
	buf.Reset()
	golog.Warnf("warning message %s", "warning")
	require.Contains(buf.String(), "warning message warning")
	buf.Reset()
	golog.Errorf("error message %s", "error")
	require.Contains(buf.String(), "error message error")
	buf.Reset()
	golog.WithValues("err", errors.New("error")).Debugf("debug message")
	require.Empty(buf.String())
	golog.WithValues("err", errors.New("error")).WithValues("c", false).Warnf("warn message")
	require.Contains(buf.String(), "warn message")
	require.Contains(buf.String(), "err=\x1b[0merror")
	require.Contains(buf.String(), "c=\x1b[0mfalse")
	buf.Reset()
	golog.WithValues("a", 1, "b", true).Infof("info message with %d fields", 2)
	require.Contains(buf.String(), "info message with 2 fields")
	require.Contains(buf.String(), "\x1b[36ma=\x1b[0m1")
	require.Contains(buf.String(), "\x1b[36mb=\x1b[0mtrue")
	buf.Reset()
	golog.Debugf("debug message")
	require.Empty(buf.String())
	l := golog.WithValues("a", 1, "b", 3)
	l.Error("error message")
	require.Contains(buf.String(), "error message")
	require.Contains(buf.String(), "\x1b[36ma=\x1b[0m1")
	require.Contains(buf.String(), "\x1b[36mb=\x1b[0m3")
	buf.Reset()
	l.WithValues("c", false).Warn("warn message")
	require.Contains(buf.String(), "warn message")
	require.Contains(buf.String(), "\x1b[36ma=\x1b[0m1")
	require.Contains(buf.String(), "\x1b[36mb=\x1b[0m3")

}

func TestDebugGlobal(t *testing.T) {
	l := golog.WithValues("a", 1, "b", 3)
	l.Error("error message")
	l.Warn("warn message", "c", false)
}

func TestDebugGlobal2(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	err := golog.LoadConfig("testdata/yaml_001.yml")
	require.NoError(err)
	for i := 0; i < 3; i++ {
		golog.Info("abcde1234")

	}
}
func TestGlobal_Panic(t *testing.T) {
	var buf bytes.Buffer
	require := require.New(t)
	cfg := golog.Config{
		Level:        golog.INFO,
		Encoding:     golog.TextEncoding,
		CallerLevels: []golog.Level{},
		TextEncoder: golog.TextEncoderConfig{
			DisableTimestamp: true,
			DisableColor:     true,
		},
		Handler: golog.HandlerConfig{
			Type:   "custom",
			Writer: &buf,
		},
	}
	log, err := golog.NewLoggerByConfig("test", cfg)
	require.NoError(err)
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		log.Panicf("panic message")
	}()
	require.NotNil(recovered)
	require.Equal("PNIC panic message\n", buf.String())
	require.Equal("panic message", recovered)
}

func TestGlobal_Fatal(t *testing.T) {
	if os.Getenv("BE_FATAL") == "1" {
		golog.Fatalf("%s", "fatal")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestGlobal_Fatal")
	cmd.Env = append(os.Environ(), "BE_FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestGlobalLogRaces(t *testing.T) {
	require := require.New(t)
	cfg := golog.Config{
		Level:    golog.INFO,
		Encoding: golog.TextEncoding,
		TextEncoder: golog.TextEncoderConfig{
			DisableTimestamp: true,
		},
		Handler: golog.HandlerConfig{
			Type: "file",
			File: golog.FileConfig{
				Path: "",
			},
		},
	}
	log, err := golog.NewLoggerByConfig("test", cfg)
	require.NoError(err)
	f := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			log.WithValues("a", 1).Info("info")
		}
	}

	var wg sync.WaitGroup
	wg.Add(4)
	go f(&wg)
	go f(&wg)
	go f(&wg)
	go f(&wg)
	wg.Wait()
}

func BenchmarkGlobal(b *testing.B) {
	defer resetConfigs()
	require := require.New(b)
	err := golog.LoadConfig("testdata/bench.yml")
	require.NoError(err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		golog.Info("abcde1234")
	}
}

func BenchmarkGlobal_WithField(b *testing.B) {
	defer resetConfigs()
	require := require.New(b)
	err := golog.LoadConfig("testdata/bench.yml")
	require.NoError(err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		golog.Info("abcde1234", "k", 1, "a", "c", "b", true)
	}
}

func BenchmarkSlog_WithValues(b *testing.B) {
	defer resetConfigs()
	olog := slog.NewTextHandler(io.Discard, nil)
	slog := slog.New(olog)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		slog.With("k", 1, "a", "c", "b", true).Info("abcde1234")
	}
}

func TestDebug(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	err := golog.LoadConfig("testdata/debug.yml")
	require.NoError(err)
	for i := 0; i < 2; i++ {
		golog.Info("abcde1233", "k", 1)
	}
	for i := 0; i < 2; i++ {
		golog.Info("abcde1234", "k", 1)
	}
}

package golog

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/millken/golog/internal/config"
	"github.com/millken/golog/internal/log"
	"github.com/stretchr/testify/require"
)

func TestGlobalLog(t *testing.T) {
	Debug("debug message")
	Info("info message")

	Warn("warning message")
	Error("error message")
	Debugf("debug message %s", "debug")
	Infof("info message %s", "info")

	Warnf("warning message %s", "warning")
	Errorf("error message %s", "error")
	WithField("err", errors.New("error")).Debugf("debug message")
	WithField("err", errors.New("error")).WithField("c", false).Warnf("warn message")
	WithFields(F("a", 1), F("b", true)).Infof("info message with %d fields", 2)
	Debugf("debug message")

	l := WithFields(F("a", 1), F("b", 3))
	l.Error("error message")
	l.WithField("c", false).Warn("warn message")
}

func TestGlobal_Panic(t *testing.T) {
	var buf bytes.Buffer
	require := require.New(t)
	cfg := config.Config{
		Level:    log.INFO,
		Encoding: "console",
		ConsoleEncoderConfig: config.ConsoleEncoderConfig{
			DisableTimestamp: true,
			DisableColor:     true,
		},
		Writer: config.WriterConfig{
			Type:         "custom",
			CustomWriter: &buf,
		},
	}
	log, err := NewLoggerByConfig("test", cfg)
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
		Fatalf("%s", "fatal")
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
	f := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			log.Info("info")
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

func BenchmarkGlobalLogger(b *testing.B) {
	require := require.New(b)
	err := LoadConfig("internal/config/testdata/bench.yml")
	require.NoError(err)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Infof("abcde1234")
	}
}

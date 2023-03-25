package golog

import (
	"bytes"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func resetConfigs() {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs = newConfigs()
}

func TestSetConfig(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	SetLevel(INFO)
	SetEncoding(TextEncoding)
	SetCallerLevels(PANIC, FATAL, ERROR, WARNING, INFO, DEBUG)
	SetStacktraceLevels(PANIC, FATAL, ERROR, WARNING)
	buf := &bytes.Buffer{}
	SetWriter(buf)
	SetConsoleEncoderConfig(TextEncoderConfig{DisableTimestamp: true})
	SetJSONEncoderConfig(JSONEncoderConfig{DisableTimestamp: true})
	require.Equal(INFO, configs.Default.Level)
	require.Equal(TextEncoding, configs.Default.Encoding)
	sort.Slice(configs.Default.CallerLevels, func(i, j int) bool {
		return configs.Default.CallerLevels[i] < configs.Default.CallerLevels[j]
	})
	require.Equal([]Level{PANIC, FATAL, ERROR, WARNING, INFO, DEBUG}, configs.Default.CallerLevels)
	require.Equal([]Level{PANIC, FATAL, ERROR, WARNING}, configs.Default.StacktraceLevels)
	require.True(configs.Default.TextEncoderConfig.DisableTimestamp)
	require.True(configs.Default.JSONEncoderConfig.DisableTimestamp)

	SetModuleConfig("mudule/1", Config{Level: DEBUG, Encoding: JSONEncoding})
	require.Equal(1, len(configs.Modules))
	require.Equal(DEBUG, configs.Modules["mudule/1"].Level)
}

func TestConfig(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	err := LoadConfig("./testdata/yaml_001.yml")
	require.NoError(err)
	require.Equal(INFO, configs.Default.Level)
	require.Equal(TextEncoding, configs.Default.Encoding)
	sort.Slice(configs.Default.CallerLevels, func(i, j int) bool {
		return configs.Default.CallerLevels[i] < configs.Default.CallerLevels[j]
	})
	require.Equal([]Level{PANIC, FATAL, ERROR, WARNING, INFO, DEBUG},
		configs.Default.CallerLevels,
	)
	sort.Slice(configs.Default.StacktraceLevels, func(i, j int) bool {
		return configs.Default.StacktraceLevels[i] < configs.Default.StacktraceLevels[j]
	})
	require.Equal([]Level{PANIC, FATAL, ERROR, WARNING},
		configs.Default.StacktraceLevels,
	)
	require.Equal("file", configs.Default.Writer.Type)
	require.Equal("/var/log/golog.log", configs.Default.Writer.FileConfig.Path)

	require.Equal(1, len(configs.Modules))
	require.Equal(DEBUG, configs.Modules["mudule/1"].Level)
	require.Equal(JSONEncoding, configs.Modules["mudule/1"].Encoding)

	cfg := GetModuleConfig("mudule/1")
	require.Equal(DEBUG, cfg.Level)
	require.Equal(JSONEncoding, cfg.Encoding)
}

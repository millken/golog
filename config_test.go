package golog

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func resetConfigs() {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	configs = newConfigs()
}

func TestConfig(t *testing.T) {
	defer resetConfigs()
	require := require.New(t)
	err := LoadConfig("./testdata/yaml_001.yml")
	require.NoError(err)
	require.Equal(INFO, configs.Default.Level)
	require.Equal("console", configs.Default.Encoding)
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
	require.Equal("json", configs.Modules["mudule/1"].Encoding)

	cfg := GetModuleConfig("mudule/1")
	require.Equal(DEBUG, cfg.Level)
	require.Equal("json", cfg.Encoding)
}

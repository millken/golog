package config

import (
	"sort"
	"testing"

	"github.com/millken/golog/log"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	require := require.New(t)
	err := Load("../testdata/yaml_001.yml")
	require.NoError(err)
	require.Equal(log.INFO, configs.Default.Level)
	require.Equal("console", configs.Default.Encoding)
	sort.Slice(configs.Default.CallerLevels, func(i, j int) bool {
		return configs.Default.CallerLevels[i] < configs.Default.CallerLevels[j]
	})
	require.Equal([]log.Level{log.PANIC, log.FATAL, log.ERROR, log.WARNING, log.INFO, log.DEBUG},
		configs.Default.CallerLevels,
	)
	sort.Slice(configs.Default.StacktraceLevels, func(i, j int) bool {
		return configs.Default.StacktraceLevels[i] < configs.Default.StacktraceLevels[j]
	})
	require.Equal([]log.Level{log.PANIC, log.FATAL, log.ERROR, log.WARNING},
		configs.Default.StacktraceLevels,
	)
	require.Equal("file", configs.Default.Writer.Type)
	require.Equal("/var/log/golog.log", configs.Default.Writer.FileConfig.Path)

	require.Equal(1, len(configs.Modules))
	require.Equal(log.DEBUG, configs.Modules["mudule/1"].Level)
	require.Equal("json", configs.Modules["mudule/1"].Encoding)

	cfg := GetModuleConfig("mudule/1")
	require.Equal(log.DEBUG, cfg.Level)
	require.Equal("json", cfg.Encoding)
}

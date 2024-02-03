package golog_test

// func TestSetConfig(t *testing.T) {
// 	defer resetConfigs()
// 	configs := golog.GetConfigs()
// 	require := require.New(t)
// 	golog.SetLevel(golog.INFO)
// 	golog.SetEncoding(golog.TextEncoding)
// 	golog.SetCallerLevels(golog.PANIC, golog.FATAL, golog.ERROR, golog.WARNING, golog.INFO, golog.DEBUG)
// 	golog.SetStacktraceLevels(golog.PANIC, golog.FATAL, golog.ERROR, golog.WARNING)
// 	buf := &bytes.Buffer{}
// 	golog.SetWriter(buf)
// 	golog.SetTextEncoderConfig(golog.TextEncoderConfig{DisableTimestamp: true})
// 	golog.SetJSONEncoderConfig(golog.JSONEncoderConfig{DisableTimestamp: true})
// 	require.Equal(golog.INFO, configs.Default.Level)
// 	require.Equal(golog.TextEncoding, configs.Default.Encoding)
// 	sort.Slice(configs.Default.CallerLevels, func(i, j int) bool {
// 		return configs.Default.CallerLevels[i] < configs.Default.CallerLevels[j]
// 	})
// 	require.Equal([]golog.Level{golog.PANIC, golog.FATAL, golog.ERROR, golog.WARNING, golog.INFO, golog.DEBUG}, configs.Default.CallerLevels)
// 	require.Equal([]golog.Level{golog.PANIC, golog.FATAL, golog.ERROR, golog.WARNING}, configs.Default.StacktraceLevels)
// 	require.True(configs.Default.TextEncoder.DisableTimestamp)
// 	require.True(configs.Default.JSONEncoder.DisableTimestamp)

// 	golog.SetModuleConfig("mudule/1", golog.Config{Level: golog.DEBUG, Encoding: golog.JSONEncoding})
// 	require.Equal(1, len(configs.Modules))
// 	require.Equal(golog.DEBUG, configs.Modules["mudule/1"].Level)
// }

// func TestConfig(t *testing.T) {
// 	defer resetConfigs()
// 	require := require.New(t)
// 	err := golog.LoadConfig("./testdata/yaml_001.yml")
// 	require.NoError(err)
// 	configs := golog.GetConfigs()
// 	require.Equal(golog.INFO, configs.Default.Level)
// 	require.Equal(golog.TextEncoding, configs.Default.Encoding)
// 	sort.Slice(configs.Default.CallerLevels, func(i, j int) bool {
// 		return configs.Default.CallerLevels[i] < configs.Default.CallerLevels[j]
// 	})
// 	require.Equal([]golog.Level{golog.PANIC, golog.FATAL, golog.ERROR, golog.WARNING, golog.INFO, golog.DEBUG},
// 		configs.Default.CallerLevels,
// 	)
// 	sort.Slice(configs.Default.StacktraceLevels, func(i, j int) bool {
// 		return configs.Default.StacktraceLevels[i] < configs.Default.StacktraceLevels[j]
// 	})
// 	require.Equal([]golog.Level{golog.PANIC, golog.FATAL, golog.ERROR, golog.WARNING},
// 		configs.Default.StacktraceLevels,
// 	)
// 	require.Equal("file", configs.Default.Handler.Type)
// 	require.Equal("/var/log/golog.log", configs.Default.Handler.File.Path)

// 	require.Equal(1, len(configs.Modules))
// 	require.Equal(golog.DEBUG, configs.Modules["mudule/1"].Level)
// 	require.Equal(golog.JSONEncoding, configs.Modules["mudule/1"].Encoding)

// 	cfg := golog.GetModuleConfig("mudule/1")
// 	require.Equal(golog.DEBUG, cfg.Level)
// 	require.Equal(golog.JSONEncoding, cfg.Encoding)
// }

// func TestConfig2(t *testing.T) {
// 	defer resetConfigs()
// 	require := require.New(t)
// 	err := golog.LoadConfig("./testdata/yaml_002.yml")
// 	require.NoError(err)
// 	configs := golog.GetConfigs()
// 	require.Equal(golog.INFO, configs.Default.Level)
// 	require.Equal(golog.JSONEncoding, configs.Default.Encoding)

// 	require.Equal("rotateFile", configs.Default.Handler.Type)
// 	require.Equal("", configs.Default.Handler.File.Path)
// 	require.Equal("/var/log/golog.log", configs.Default.Handler.RotateFile.Filename)
// 	require.Equal(3, configs.Default.Handler.RotateFile.MaxBackups)
// 	require.Equal("2006-01-02", configs.Default.Handler.RotateFile.BackupTimeFormat)
// 	require.True(configs.Default.Handler.RotateFile.LocalTime)
// 	require.True(configs.Default.Handler.RotateFile.Async)
// 	require.Equal(1, len(configs.Modules))
// 	require.Equal(golog.DEBUG, configs.Modules["mudule/1"].Level)
// 	require.Equal(golog.JSONEncoding, configs.Modules["mudule/1"].Encoding)

// 	cfg := golog.GetModuleConfig("mudule/1")
// 	require.Equal(golog.DEBUG, cfg.Level)
// 	require.Equal(golog.JSONEncoding, cfg.Encoding)
// }

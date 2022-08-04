package meta

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/millken/golog/internal/log"
	"github.com/stretchr/testify/require"
)

func TestLevels(t *testing.T) {
	module := "sample-module-critical"
	SetLevel(module, log.FATAL)
	require.Equal(t, log.FATAL, GetLevel(module))
	verifyLevels(t, module, []log.Level{log.FATAL}, []log.Level{log.ERROR, log.WARNING, log.INFO, log.DEBUG})

	module = "sample-module-error"
	SetLevel(module, log.ERROR)
	require.Equal(t, log.ERROR, GetLevel(module))
	verifyLevels(t, module, []log.Level{log.FATAL, log.ERROR}, []log.Level{log.WARNING, log.INFO, log.DEBUG})

	module = "sample-module-warning"
	SetLevel(module, log.WARNING)
	require.Equal(t, log.WARNING, GetLevel(module))
	verifyLevels(t, module, []log.Level{log.FATAL, log.ERROR, log.WARNING}, []log.Level{log.INFO, log.DEBUG})

	module = "sample-module-info"
	SetLevel(module, log.INFO)
	require.Equal(t, log.INFO, GetLevel(module))
	verifyLevels(t, module, []log.Level{log.FATAL, log.ERROR, log.WARNING, log.INFO}, []log.Level{log.DEBUG})

	module = "sample-module-debug"
	SetLevel(module, log.DEBUG)
	require.Equal(t, log.DEBUG, GetLevel(module))
	verifyLevels(t, module, []log.Level{log.FATAL, log.ERROR, log.WARNING, log.INFO, log.DEBUG}, []log.Level{})

	module = "def-module-debug"
	SetDefaultLevel(log.DEBUG)
	require.Equal(t, log.DEBUG, GetLevel(module))
	verifyLevels(t, module, []log.Level{log.FATAL, log.ERROR, log.WARNING, log.INFO, log.DEBUG}, []log.Level{})
}

// func TestHandlers(t *testing.T) {
// 	module := "sample-module-handler"
// 	h1 := handler.NewNull()
// 	SetHandler(module, h1)
// 	require.Equal(t, h1, GetHandler(module))

// 	h2, _ := handler.NewFile(handler.FileConfig{Path: ""})
// 	SetDefaultHandler(h2)

// 	module = "sample-module-handler2"
// 	require.Equal(t, h2, GetHandler(module))
// }

func TestCallerInfos(t *testing.T) {
	module := fmt.Sprintf("sample-module-caller-info-%d-%d", rand.Intn(1000), rand.Intn(1000)) //nolint:gosec

	require.True(t, IsCallerInfoEnabled(module, log.FATAL))
	require.True(t, IsCallerInfoEnabled(module, log.DEBUG))
	require.True(t, IsCallerInfoEnabled(module, log.INFO))
	require.True(t, IsCallerInfoEnabled(module, log.ERROR))
	require.True(t, IsCallerInfoEnabled(module, log.WARNING))

	ShowCallerInfo(module, log.FATAL)
	ShowCallerInfo(module, log.DEBUG)
	HideCallerInfo(module, log.INFO)
	HideCallerInfo(module, log.ERROR)
	HideCallerInfo(module, log.WARNING)

	require.True(t, IsCallerInfoEnabled(module, log.FATAL))
	require.True(t, IsCallerInfoEnabled(module, log.DEBUG))
	require.False(t, IsCallerInfoEnabled(module, log.INFO))
	require.False(t, IsCallerInfoEnabled(module, log.ERROR))
	require.False(t, IsCallerInfoEnabled(module, log.WARNING))

	require.True(t, IsCallerInfoEnabled(module, log.FATAL))
	require.True(t, IsCallerInfoEnabled(module, log.DEBUG))
	require.False(t, IsCallerInfoEnabled(module, log.INFO))
	require.False(t, IsCallerInfoEnabled(module, log.ERROR))
	require.False(t, IsCallerInfoEnabled(module, log.WARNING))
}

func TestParseLevel(t *testing.T) {
	verifyLevelsNoError := func(expected log.Level, levels ...string) {
		for _, level := range levels {
			actual, err := ParseLevel(level)
			require.NoError(t, err, "not supposed to fail while parsing level string [%s]", level)
			require.Equal(t, expected, actual)
		}
	}

	verifyLevelsNoError(log.FATAL, "fatal", "FATAL", "FaTaL")
	verifyLevelsNoError(log.ERROR, "error", "ERROR", "ErroR")
	verifyLevelsNoError(log.WARNING, "warning", "WARNING", "WarninG")
	verifyLevelsNoError(log.DEBUG, "debug", "DEBUG", "DebUg")
	verifyLevelsNoError(log.INFO, "info", "INFO", "iNFo")
}

func TestParseLevelError(t *testing.T) {
	verifyLevelError := func(levels ...string) {
		for _, level := range levels {
			_, err := ParseLevel(level)
			require.Error(t, err, "not supposed to succeed while parsing level string [%s]", level)
		}
	}

	verifyLevelError("", "D", "DE BUG", ".")
}

func TestParseString(t *testing.T) {
	require.Equal(t, "FATAL", ParseString(log.FATAL))
	require.Equal(t, "ERROR", ParseString(log.ERROR))
	require.Equal(t, "WARNING", ParseString(log.WARNING))
	require.Equal(t, "DEBUG", ParseString(log.DEBUG))
	require.Equal(t, "INFO", ParseString(log.INFO))
}

func verifyLevels(t *testing.T, module string, enabled, disabled []log.Level) {
	for _, level := range enabled {
		actual := IsEnabledFor(module, level)
		require.True(t, actual, "expected level [%s] to be enabled for module [%s]", ParseString(level), module)
	}

	for _, level := range disabled {
		actual := IsEnabledFor(module, level)
		require.False(t, actual, "expected level [%s] to be disabled for module [%s]", ParseString(level), module)
	}
}

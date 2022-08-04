package meta

import (
	"testing"

	"github.com/millken/golog/internal/log"
	"github.com/stretchr/testify/require"
)

func TestCallerInfoSetting(t *testing.T) {
	ci := newCallerInfo()
	mod := "sample-module-name"

	// By default caller info should be enabled if not set
	require.True(t, ci.IsCallerInfoEnabled(mod, log.DEBUG), "Callerinfo supposed to be enabled for this level")
	require.True(t, ci.IsCallerInfoEnabled(mod, log.INFO), "Callerinfo supposed to be enabled for this level")
	require.True(t, ci.IsCallerInfoEnabled(mod, log.WARNING), "Callerinfo supposed to be enabled for this level")
	require.True(t, ci.IsCallerInfoEnabled(mod, log.ERROR), "Callerinfo supposed to be enabled for this level")
	require.True(t, ci.IsCallerInfoEnabled(mod, log.FATAL), "Callerinfo supposed to be enabled for this level")
	require.True(t, ci.IsCallerInfoEnabled(mod, log.PANIC), "Callerinfo supposed to be enabled for this level")

	ci.HideCallerInfo(mod, log.DEBUG)
	require.False(t, ci.IsCallerInfoEnabled(mod, log.DEBUG), "Callerinfo supposed to be disabled for this level")

	ci.ShowCallerInfo(mod, log.DEBUG)
	require.True(t, ci.IsCallerInfoEnabled(mod, log.DEBUG), "Callerinfo supposed to be enabled for this level")

	ci.HideCallerInfo(mod, log.WARNING)
	require.False(t, ci.IsCallerInfoEnabled(mod, log.WARNING), "Callerinfo supposed to be disabled for this level")

	ci.ShowCallerInfo(mod, log.WARNING)
	require.True(t, ci.IsCallerInfoEnabled(mod, log.WARNING), "Callerinfo supposed to be enabled for this level")

	ci.HideCallerInfo(mod, log.DEBUG)
	require.False(t, ci.IsCallerInfoEnabled(mod, log.DEBUG), "Callerinfo supposed to be disabled for this level")

	ci.ShowCallerInfo(mod, log.DEBUG)
	require.True(t, ci.IsCallerInfoEnabled(mod, log.DEBUG), "Callerinfo supposed to be enabled for this level")

	// By default caller info should be enabled for any module name not set before
	moduleNames := []string{"sample-module-name-doesnt-exists", "", "@$#@$@"}
	for _, moduleName := range moduleNames {
		require.True(t, ci.IsCallerInfoEnabled(moduleName, log.INFO), "Callerinfo supposed to be enabled for this level")
		require.True(t, ci.IsCallerInfoEnabled(moduleName, log.WARNING), "Callerinfo supposed to be enabled for this level")
		require.True(t, ci.IsCallerInfoEnabled(moduleName, log.ERROR), "Callerinfo supposed to be enabled for this level")
		require.True(t, ci.IsCallerInfoEnabled(moduleName, log.FATAL), "Callerinfo supposed to be enabled for this level")
		require.True(t, ci.IsCallerInfoEnabled(moduleName, log.DEBUG), "Callerinfo supposed to be enabled for this level")
		require.True(t, ci.IsCallerInfoEnabled(moduleName, log.PANIC), "Callerinfo supposed to be enabled for this level")
	}
}

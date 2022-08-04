package meta

import (
	"testing"

	"github.com/millken/golog/internal/log"
	"github.com/stretchr/testify/require"
)

func TestLogLevels(t *testing.T) {
	mlevel := newModuledLevels()
	mlevel.SetDefaultLevel(log.INFO)

	mlevel.SetLevel("module-xyz-info", log.INFO)
	mlevel.SetLevel("module-xyz-debug", log.DEBUG)
	mlevel.SetLevel("module-xyz-error", log.ERROR)
	mlevel.SetLevel("module-xyz-warning", log.WARNING)
	mlevel.SetLevel("module-xyz-fatal", log.FATAL)

	// Run info level checks
	require.True(t, mlevel.IsEnabledFor("module-xyz-info", log.FATAL))
	require.True(t, mlevel.IsEnabledFor("module-xyz-info", log.ERROR))
	require.True(t, mlevel.IsEnabledFor("module-xyz-info", log.WARNING))
	require.True(t, mlevel.IsEnabledFor("module-xyz-info", log.INFO))
	require.False(t, mlevel.IsEnabledFor("module-xyz-info", log.DEBUG))

	// Run debug level checks
	require.True(t, mlevel.IsEnabledFor("module-xyz-debug", log.FATAL))
	require.True(t, mlevel.IsEnabledFor("module-xyz-debug", log.ERROR))
	require.True(t, mlevel.IsEnabledFor("module-xyz-debug", log.WARNING))
	require.True(t, mlevel.IsEnabledFor("module-xyz-debug", log.INFO))
	require.True(t, mlevel.IsEnabledFor("module-xyz-debug", log.DEBUG))

	// Run warning level checks
	require.True(t, mlevel.IsEnabledFor("module-xyz-warning", log.FATAL))
	require.True(t, mlevel.IsEnabledFor("module-xyz-warning", log.ERROR))
	require.True(t, mlevel.IsEnabledFor("module-xyz-warning", log.WARNING))
	require.False(t, mlevel.IsEnabledFor("module-xyz-warning", log.INFO))
	require.False(t, mlevel.IsEnabledFor("module-xyz-warning", log.DEBUG))

	// Run error level checks
	require.True(t, mlevel.IsEnabledFor("module-xyz-error", log.FATAL))
	require.True(t, mlevel.IsEnabledFor("module-xyz-error", log.ERROR))
	require.False(t, mlevel.IsEnabledFor("module-xyz-error", log.WARNING))
	require.False(t, mlevel.IsEnabledFor("module-xyz-error", log.INFO))
	require.False(t, mlevel.IsEnabledFor("module-xyz-error", log.DEBUG))

	// Run error fatal checks
	require.True(t, mlevel.IsEnabledFor("module-xyz-fatal", log.FATAL))
	require.False(t, mlevel.IsEnabledFor("module-xyz-fatal", log.ERROR))
	require.False(t, mlevel.IsEnabledFor("module-xyz-fatal", log.WARNING))
	require.False(t, mlevel.IsEnabledFor("module-xyz-fatal", log.INFO))
	require.False(t, mlevel.IsEnabledFor("module-xyz-fatal", log.DEBUG))

	// Run default log level check --> which is info level
	require.True(t, mlevel.IsEnabledFor("module-xyz-random-module", log.FATAL))
	require.True(t, mlevel.IsEnabledFor("module-xyz-random-module", log.ERROR))
	require.True(t, mlevel.IsEnabledFor("module-xyz-random-module", log.WARNING))
	require.True(t, mlevel.IsEnabledFor("module-xyz-random-module", log.INFO))
	require.False(t, mlevel.IsEnabledFor("module-xyz-random-module", log.DEBUG))
}

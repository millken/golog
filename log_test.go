package golog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLog(t *testing.T) {
	require := require.New(t)
	require.NoError(LoadConfig("./internal/config/testdata/sample.yml"))
	Debug("debug")
	Info("info")
	Warn("warn")
	Error("error")
}

func BenchmarkLog(b *testing.B) {
	require := require.New(b)
	require.NoError(LoadConfig("./internal/config/testdata/bench.yml"))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Info("info")
	}
}

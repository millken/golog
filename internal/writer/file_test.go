package writer

import (
	"bytes"
	"testing"

	"github.com/millken/golog/config"
	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {
	require := require.New(t)
	var b bytes.Buffer
	f, err := NewFile(config.FileConfig{Path: "stdout"})
	require.NoError(err)
	f.writer = &b
	_, err = f.Write([]byte("test"))
	require.NoError(err)
	require.Contains(b.String(), "test")
}

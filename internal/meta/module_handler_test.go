package meta

import (
	"testing"

	"github.com/millken/golog/internal/handler"
	"github.com/stretchr/testify/require"
)

func TestLogHandlers(t *testing.T) {
	mhandler := newModuledHandlers()
	mhandler.SetDefaultHandler(defaultHandler)
	mhandler.SetHandler("test", handler.NewNull())
	h := mhandler.GetHandler("test")
	require.NotNil(t, h)
	require.IsType(t, &handler.Null{}, h)

	h = mhandler.GetHandler("test2")
	require.Equal(t, defaultHandler, h)
}

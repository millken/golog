package writer

import (
	"bytes"
	"testing"

	"github.com/millken/golog"
	"github.com/stretchr/testify/assert"
)

func TestDifferentLevelsGoToDifferentWriters(t *testing.T) {
	var a, b bytes.Buffer

	log := golog.NewLogger()
	hand1 := &Handler{
		Output: &a,
	}
	hand1.SetLevels(golog.WarnLevel)
	hand1.SetFormatter(&golog.TextFormatter{
		DisableTimestamp: true,
		NoColor:          true,
	})

	log.AddHandler(hand1)

	hand2 := &Handler{
		Output: &b,
	}
	hand2.SetLevels(golog.InfoLevel)
	hand2.SetFormatter(&golog.TextFormatter{
		DisableTimestamp: true,
		NoColor:          true,
	})
	log.AddHandler(hand2)
	log.Warnf("send to a")
	log.Infof("send to b")

	assert.Equal(t, a.String(), "warn send to a\n")
	assert.Equal(t, b.String(), "info send to b\n")
}

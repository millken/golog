package golog

import (
	"io"
	"testing"
)

func TestStdLog(t *testing.T) {
	Debug("std debug message")
	Info("std info message")
	Warn("std warning message")
	Error("std error message")
	SetLevel(DebugLevel)
	Debug("std debug message")
	Info("std debug message with 2 fields", Field("a", 1), Field("b", true))
	SetOutput(io.Discard)
	Debug("std debug message")
}

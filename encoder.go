//go:build !binary_log
// +build !binary_log

package golog

import (
	"sync/atomic"

	"github.com/millken/x/buffer"
)

// Encoder is a interface for encoding log entry.

type Encoder interface {
	Encode(*buffer.Buffer, Record) error
}

var (
	defaultEncoderType atomic.Bool
	defaultEncoderText atomic.Pointer[EncoderText]
	defaultEncoderJson atomic.Pointer[EncoderJSON]
)

func init() {
	defaultEncoderType.Store(false)
	defaultEncoderText.Store(NewEncoderText(EncoderTextConfig{}))
}

func defaultEncoder() Encoder {
	if defaultEncoderType.Load() {
		return defaultEncoderJson.Load()
	}
	return defaultEncoderText.Load()
}

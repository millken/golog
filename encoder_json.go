//go:build !binary_log
// +build !binary_log

package golog

// encoder_json.go file contains bindings to generate
// JSON encoded byte stream.

import (
	"github.com/millken/golog/internal/json"
)

var (
	_ encoder = (*json.Encoder)(nil)

	enc = json.Encoder{}
)

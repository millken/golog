package golog

import (
	"unsafe"
)

func B2S(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// S2B converts string to a byte slice without memory allocation.
func S2B(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

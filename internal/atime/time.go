package atime

import (
	"sync/atomic"
	"time"

	"github.com/millken/gosync"
)

var (
	_time  int64
	_cache = new(gosync.Ontime)
)

func Now(native bool) time.Time {
	if native {
		return time.Now()
	}
	_cache.Do(700*time.Millisecond, func() {
		atomic.StoreInt64(&_time, time.Now().UnixNano())
	})
	return time.Unix(0, atomic.LoadInt64(&_time))
}

package golog

import (
	"sync"
	"sync/atomic"
)

// Encoding is the encoding type.
type Encoding string

const (
	defaultLogLevel = INFO
	defaultModule   = "-"
	//JSONEncoding is the json encoding.
	JSONEncoding Encoding = "json"
	//TextEncoding is the text encoding.
	TextEncoding Encoding = "text"
)

var (
	rwmutex            = &sync.RWMutex{}
	enableNativeTime   bool
	defaultConfigValue atomic.Value
)
var (
	config = Config{
		Shortfile: true,
		Longfile:  true,
		Stack:     true,
		Datetime:  true,
		Timestamp: true,
		UTC:       true,
		Function:  true,
		Level:     defaultLogLevel,
	}
)

func defaultConfig() Config {
	return defaultConfigValue.Load().(Config)
}

type Config struct {
	Fields    []Field
	Datetime  bool
	Timestamp bool
	UTC       bool
	Shortfile bool
	Longfile  bool
	Stack     bool
	Function  bool
	Level     Level

	calldepth int
}

package golog

import (
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		// var b bytes.Buffer
		// b.Grow(64)
		return &bytes.Buffer{}
	},
}

var mapPool = &sync.Pool{
	New: func() interface{} {
		return make(map[string]interface{})
	},
}

type JSONFormatter struct {
	// EnableCaller enabled caller
	EnableCaller bool
}

const levelName = "entry.Level"

func (f *JSONFormatter) Format(entry *Entry) error {
	entry.Data = enc.AppendBeginMarker(entry.Data)
	entry.Data = appendKeyVal(entry.Data, TimestampFieldName, &entry.Timestamp)
	entry.Data = appendKeyVal(entry.Data, LevelFieldName, entry.Level.String())
	entry.Data = appendKeyVal(entry.Data, MessageFieldName, &entry.Message)

	if f.EnableCaller {
		file, line := entry.GetCaller(CallerSkipFrameCount)
		c := file + ":" + strconv.Itoa(line)
		if len(c) > 0 {
			if cwd, err := os.Getwd(); err == nil {
				if rel, err := filepath.Rel(cwd, c); err == nil {
					c = rel
				}
			}
		}
		entry.Data = appendKeyVal(entry.Data, CallerFieldName, c)
	}
	entry.Data = appendFields(entry.Data, entry.Fields[:entry.fieldsLen]...)
	entry.Data = enc.AppendEndMarker(entry.Data)
	entry.Data = enc.AppendLineBreak(entry.Data)
	return nil
}

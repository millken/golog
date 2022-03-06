package golog

import (
	"path/filepath"
	"strconv"
)

type JSONFormatter struct {
	// EnableCaller enabled caller
	EnableCaller     bool
	DisableTimestamp bool
}

func (f *JSONFormatter) Format(entry *Entry) error {
	entry.Data = enc.AppendBeginMarker(entry.Data)
	if !f.DisableTimestamp {
		entry.Data = appendKeyVal(entry.Data, TimestampFieldName, &entry.Timestamp)
	}
	entry.Data = appendKeyVal(entry.Data, LevelFieldName, entry.Level.String())
	entry.Data = appendKeyVal(entry.Data, MessageFieldName, &entry.Message)

	if f.EnableCaller {
		file, line := caller(CallerSkipFrameCount)
		c := file + ":" + strconv.Itoa(line)
		if len(c) > 0 {
			if rel, err := filepath.Rel(_cwd, c); err == nil {
				c = rel
			}
		}
		entry.Data = appendKeyVal(entry.Data, CallerFieldName, c)
	}
	entry.Data = appendFields(entry.Data, entry.Fields[:entry.fieldsLen]...)
	entry.Data = enc.AppendEndMarker(entry.Data)
	entry.Data = enc.AppendLineBreak(entry.Data)
	return nil
}

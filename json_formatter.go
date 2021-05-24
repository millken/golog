package golog

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"unsafe"

	"github.com/mailru/easyjson"
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

func bytesToString(bytes []byte) (s string) {
	slice := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	str := (*reflect.StringHeader)(unsafe.Pointer(&s))
	str.Data = slice.Data
	str.Len = slice.Len
	runtime.KeepAlive(&bytes) // this line is essential.
	return s
}

func (f *JSONFormatter) Format(entry *Entry) error {
	buff := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buff)
	buff.Reset()

	mapLen := len(entry.Fields) + 3
	if f.EnableCaller {
		mapLen = +1
	}
	fields := make(Fields, mapLen)

	for _, field := range entry.Fields {
		fields[b2s(field.key)] = field.value
	}

	fields[MessageFieldName] = bytesToString(entry.Data)
	fields[LevelFieldName] = entry.Level.String()
	fields[TimestampFieldName] = entry.Timestamp

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
		fields[CallerFieldName] = c

	}
	data, err := easyjson.Marshal(fields)

	if err != nil {
		return err
	}
	data = append(data, []byte("\n")...)
	_, err = entry.Write(data)
	return err
}

//go:build !binary_log
// +build !binary_log

package golog

import (
	"net"
	"time"
	"unsafe"

	"github.com/millken/golog/internal/json"
)

var (
	_ encoder = (*json.Encoder)(nil)

	defaultCallerSkip int = 4

	enc = json.Encoder{}

	// ErrorMarshalFunc allows customization of global error marshaling
	ErrorMarshalFunc = func(err error) interface{} {
		return err
	}
)

type encoder interface {
	AppendArrayDelim(dst []byte) []byte
	AppendArrayEnd(dst []byte) []byte
	AppendArrayStart(dst []byte) []byte
	AppendBeginMarker(dst []byte) []byte
	AppendBool(dst []byte, val bool) []byte
	AppendBools(dst []byte, vals []bool) []byte
	AppendBytes(dst, s []byte) []byte
	AppendDuration(dst []byte, d time.Duration, unit time.Duration, useInt bool) []byte
	AppendDurations(dst []byte, vals []time.Duration, unit time.Duration, useInt bool) []byte
	AppendEndMarker(dst []byte) []byte
	AppendFloat32(dst []byte, val float32) []byte
	AppendFloat64(dst []byte, val float64) []byte
	AppendFloats32(dst []byte, vals []float32) []byte
	AppendFloats64(dst []byte, vals []float64) []byte
	AppendHex(dst, s []byte) []byte
	AppendIPAddr(dst []byte, ip net.IP) []byte
	AppendIPPrefix(dst []byte, pfx net.IPNet) []byte
	AppendInt(dst []byte, val int) []byte
	AppendInt16(dst []byte, val int16) []byte
	AppendInt32(dst []byte, val int32) []byte
	AppendInt64(dst []byte, val int64) []byte
	AppendInt8(dst []byte, val int8) []byte
	AppendInterface(dst []byte, i interface{}) []byte
	AppendInts(dst []byte, vals []int) []byte
	AppendInts16(dst []byte, vals []int16) []byte
	AppendInts32(dst []byte, vals []int32) []byte
	AppendInts64(dst []byte, vals []int64) []byte
	AppendInts8(dst []byte, vals []int8) []byte
	AppendKey(dst []byte, key string) []byte
	AppendLineBreak(dst []byte) []byte
	AppendMACAddr(dst []byte, ha net.HardwareAddr) []byte
	AppendNil(dst []byte) []byte
	AppendObjectData(dst []byte, o []byte) []byte
	AppendString(dst []byte, s string) []byte
	AppendStrings(dst []byte, vals []string) []byte
	AppendTime(dst []byte, t time.Time, format string) []byte
	AppendTimes(dst []byte, vals []time.Time, format string) []byte
	AppendUint(dst []byte, val uint) []byte
	AppendUint16(dst []byte, val uint16) []byte
	AppendUint32(dst []byte, val uint32) []byte
	AppendUint64(dst []byte, val uint64) []byte
	AppendUint8(dst []byte, val uint8) []byte
	AppendUints(dst []byte, vals []uint) []byte
	AppendUints16(dst []byte, vals []uint16) []byte
	AppendUints32(dst []byte, vals []uint32) []byte
	AppendUints64(dst []byte, vals []uint64) []byte
	AppendUints8(dst []byte, vals []uint8) []byte
}

func isNilValue(i interface{}) bool {
	return (*[2]uintptr)(unsafe.Pointer(&i))[1] == 0
}

func appendVal(dst []byte, value interface{}) []byte {
	switch val := value.(type) {
	case string:
		dst = enc.AppendString(dst, val)
	case []byte:
		dst = enc.AppendBytes(dst, val)
	case error:
		switch m := ErrorMarshalFunc(val).(type) {
		case error:
			if m == nil || isNilValue(m) {
				dst = enc.AppendNil(dst)
			} else {
				dst = enc.AppendString(dst, m.Error())
			}
		case string:
			dst = enc.AppendString(dst, m)
		default:
			dst = enc.AppendInterface(dst, m)
		}
	case []error:
		dst = enc.AppendArrayStart(dst)
		for i, err := range val {
			switch m := ErrorMarshalFunc(err).(type) {
			case error:
				if m == nil || isNilValue(m) {
					dst = enc.AppendNil(dst)
				} else {
					dst = enc.AppendString(dst, m.Error())
				}
			case string:
				dst = enc.AppendString(dst, m)
			default:
				dst = enc.AppendInterface(dst, m)
			}

			if i < (len(val) - 1) {
				enc.AppendArrayDelim(dst)
			}
		}
		dst = enc.AppendArrayEnd(dst)
	case bool:
		dst = enc.AppendBool(dst, val)
	case int:
		dst = enc.AppendInt(dst, val)
	case int8:
		dst = enc.AppendInt8(dst, val)
	case int16:
		dst = enc.AppendInt16(dst, val)
	case int32:
		dst = enc.AppendInt32(dst, val)
	case int64:
		dst = enc.AppendInt64(dst, val)
	case uint:
		dst = enc.AppendUint(dst, val)
	case uint8:
		dst = enc.AppendUint8(dst, val)
	case uint16:
		dst = enc.AppendUint16(dst, val)
	case uint32:
		dst = enc.AppendUint32(dst, val)
	case uint64:
		dst = enc.AppendUint64(dst, val)
	case float32:
		dst = enc.AppendFloat32(dst, val)
	case float64:
		dst = enc.AppendFloat64(dst, val)
	case time.Time:
		dst = enc.AppendTime(dst, val, TimeFieldFormat)
	case time.Duration:
		dst = enc.AppendDuration(dst, val, time.Millisecond, false)
	case *string:
		if val != nil {
			dst = enc.AppendString(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *bool:
		if val != nil {
			dst = enc.AppendBool(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *int:
		if val != nil {
			dst = enc.AppendInt(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *int8:
		if val != nil {
			dst = enc.AppendInt8(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *int16:
		if val != nil {
			dst = enc.AppendInt16(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *int32:
		if val != nil {
			dst = enc.AppendInt32(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *int64:
		if val != nil {
			dst = enc.AppendInt64(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *uint:
		if val != nil {
			dst = enc.AppendUint(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *uint8:
		if val != nil {
			dst = enc.AppendUint8(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *uint16:
		if val != nil {
			dst = enc.AppendUint16(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *uint32:
		if val != nil {
			dst = enc.AppendUint32(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *uint64:
		if val != nil {
			dst = enc.AppendUint64(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *float32:
		if val != nil {
			dst = enc.AppendFloat32(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *float64:
		if val != nil {
			dst = enc.AppendFloat64(dst, *val)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *time.Time:
		if val != nil {
			dst = enc.AppendTime(dst, *val, TimeFieldFormat)
		} else {
			dst = enc.AppendNil(dst)
		}
	case *time.Duration:
		if val != nil {
			dst = enc.AppendDuration(dst, *val, time.Millisecond, false)
		} else {
			dst = enc.AppendNil(dst)
		}
	case []string:
		dst = enc.AppendStrings(dst, val)
	case []bool:
		dst = enc.AppendBools(dst, val)
	case []int:
		dst = enc.AppendInts(dst, val)
	case []int8:
		dst = enc.AppendInts8(dst, val)
	case []int16:
		dst = enc.AppendInts16(dst, val)
	case []int32:
		dst = enc.AppendInts32(dst, val)
	case []int64:
		dst = enc.AppendInts64(dst, val)
	case []uint:
		dst = enc.AppendUints(dst, val)
	// case []uint8:
	// 	dst = enc.AppendUints8(dst, val)
	case []uint16:
		dst = enc.AppendUints16(dst, val)
	case []uint32:
		dst = enc.AppendUints32(dst, val)
	case []uint64:
		dst = enc.AppendUints64(dst, val)
	case []float32:
		dst = enc.AppendFloats32(dst, val)
	case []float64:
		dst = enc.AppendFloats64(dst, val)
	case []time.Time:
		dst = enc.AppendTimes(dst, val, TimeFieldFormat)
	case []time.Duration:
		dst = enc.AppendDurations(dst, val, time.Millisecond, false)
	case nil:
		dst = enc.AppendNil(dst)
	case net.IP:
		dst = enc.AppendIPAddr(dst, val)
	case net.IPNet:
		dst = enc.AppendIPPrefix(dst, val)
	case net.HardwareAddr:
		dst = enc.AppendMACAddr(dst, val)
	default:
		dst = enc.AppendInterface(dst, val)
	}
	return dst
}

func appendKeyVal(dst []byte, key string, value interface{}) []byte {
	dst = enc.AppendKey(dst, key)
	dst = appendVal(dst, value)

	return dst
}

func appendFields(dst []byte, fields ...Field) []byte {
	for _, field := range fields {
		dst = appendKeyVal(dst, field.Key, field.Val)
	}
	return dst
}

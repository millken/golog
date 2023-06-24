package fasttime

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestUnixTimestamp(t *testing.T) {
	tsExpected := uint64(time.Now().Unix())
	ts := UnixTimestamp()
	if ts-tsExpected > 1 {
		t.Fatalf("unexpected UnixTimestamp; got %d; want %d", ts, tsExpected)
	}
}

func TestUnixDate(t *testing.T) {
	dateExpected := uint64(time.Now().Unix() / (24 * 3600))
	date := UnixDate()
	if date-dateExpected > 1 {
		t.Fatalf("unexpected UnixDate; got %d; want %d", date, dateExpected)
	}
}

func TestUnixHour(t *testing.T) {
	hourExpected := uint64(time.Now().Unix() / 3600)
	hour := UnixHour()
	if hour-hourExpected > 1 {
		t.Fatalf("unexpected UnixHour; got %d; want %d", hour, hourExpected)
	}
}

func TestNow(t *testing.T) {
	nowExpected := time.Now()
	now := Now()
	if now.Sub(nowExpected) > time.Second {
		t.Fatalf("unexpected Now; got %v; want %v", now, nowExpected)
	}
}
func BenchmarkUnixTimestamp(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var ts uint64
		for pb.Next() {
			ts += UnixTimestamp()
		}
		atomic.StoreUint64(&Sink, ts)
	})
}

func BenchmarkTimeNowUnix(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		var ts uint64
		for pb.Next() {
			ts += uint64(time.Now().Unix())
		}
		atomic.StoreUint64(&Sink, ts)
	})
}

// Sink should prevent from code elimination by optimizing compiler
var Sink uint64

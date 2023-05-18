package atime

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	t1 := Now(false)
	t2 := time.Now()
	if t1.Unix() != t2.Unix() {
		t.Fatal("not equal", t1, t2)
	}
}

func BenchmarkNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Now(false)
	}
}

package cache

import (
	"errors"
	"testing"
	"time"
)

func TestCachePutGet(t *testing.T) {
	c := NewSingleCache()

	c.Put("a", "b")
	v, err := c.Get("a")
	if err != nil {
		t.Error("cache entry missing", err)
	}

	if v != "b" {
		t.Error("cache entry missing", err)
	}
}

func TestCachePutWithExp(t *testing.T) {
	c := NewSingleCache()

	k := "a"
	v := "b"
	dur := time.Second
	c.PutWithExp(k, v, dur)
	time.Sleep(time.Second)

	_, err := c.Get(k)
	if !errors.Is(err, ErrNoEntry) {
		t.Error("cache expiry failed")
	}
}

func BenchmarkCachePut(b *testing.B) {
	b.ReportAllocs()
	c := NewSingleCache()

	for i := range b.N {
		c.Put(i, i)
	}
}

func BenchmarkCachePutWithExp(b *testing.B) {
	b.ReportAllocs()
	c := NewSingleCache()

	for i := range b.N {
		c.PutWithExp(i, i, time.Second)
	}
}

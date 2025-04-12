package cache

import (
	"container/list"
	"errors"
	"fmt"
	"testing"
)

func TestLRUCachePut(t *testing.T) {
	c := NewLRU[int, int](100)
	c.Put(1, 1)
	c.Put(2, 2)
	v, err := c.Get(1)

	if errors.Is(err, ErrNoEntry) {
		t.Error("key value missing")
	}

	if v != 1 {
		t.Error("value not matching")
	}

	v, err = c.Get(2)

	if errors.Is(err, ErrNoEntry) {
		t.Error("key value missing")
	}

	if v != 2 {
		t.Error("value not matching")
	}

	v, err = c.Get(3)

	if !errors.Is(err, ErrNoEntry) {
		t.Error("Key should not be present")
	}
}

func TestLRUCachePutMany(t *testing.T) {
	c := NewLRU[int, int](3)
	for i := range 1 {
		c.Put(i, i)
	}

	if !c.WithinSize() {
		t.Error("cache size is not within limit")
	}

	// for i := range 1001 {
	// 	v, err := c.Get(i)
	//
	// 	if i >= 0 && i <= 900 {
	// 		if !errors.Is(err, ErrNoEntry) {
	// 			t.Error("Key should not be present")
	// 		}
	// 	} else {
	// 		if v != i {
	// 			t.Error("value not matching")
	// 		}
	// 	}
	// }
}

func TestList(t *testing.T) {
	l := list.New()
	fmt.Println(l.Len())

	l.PushFront(1)
	fmt.Println(l.Back())

	l.PushFront(2)
	fmt.Println(l.Back())

	l.PushFront(3)
	fmt.Println(l.Back())

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e, e.Value)
	}
}

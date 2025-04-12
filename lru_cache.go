package cache

import (
	"container/list"
	"errors"
	"fmt"
	"sync"
	"time"
)

type LRUCache[K comparable, V any] struct {
	cache map[K]*Value[V]
	size  int
	list  list.List
	addr  map[K]*list.Element
	mu    sync.RWMutex
}

func NewLRU[K comparable, V any](size int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		size:  size,
		list:  *list.New(),
		addr:  make(map[K]*list.Element),
		cache: make(map[K]*Value[V]),
		mu:    sync.RWMutex{},
	}
}

func (l *LRUCache[K, V]) Put(key K, value V) {
	l.PutWithExp(key, value, 0*time.Second)
}

func (l *LRUCache[K, V]) PutWithExp(key K, value V, dur time.Duration) {
	expiration := time.Now().Add(dur)
	if dur == 0 {
		expiration = time.Time{}
	}

	_, err := l.Get(key)
	if errors.Is(err, ErrNoEntry) {
		if len(l.cache) == l.size {
			l.Evict()
		}
	}

	l.cache[key] = &Value[V]{
		v:          value,
		expiration: expiration,
	}

	l.PrintList()
	l.list.PushFront(key)
	l.PrintList()
	l.addr[key] = l.list.Front()

	return
}

func (l *LRUCache[K, V]) Evict() {
	el := l.list.Back()
	delete(l.cache, el.Value.(K))
	delete(l.addr, el.Value.(K))
	l.list.Remove(el)
}

func (l *LRUCache[K, V]) Get(key K) (value V, err error) {
	v, ok := l.cache[key]
	if !ok {
		return value, ErrNoEntry
	}

	if v.Expired() {
		el := l.addr[key]
		delete(l.cache, el.Value.(K))
		delete(l.addr, el.Value.(K))
		l.list.Remove(el)
		return value, ErrNoEntry
	}

	element := l.addr[key]
	l.list.MoveToFront(element)
	fmt.Println(l.list.Len())
	return v.v, nil
}

func (l *LRUCache[K, V]) WithinSize() bool {
	return len(l.cache) <= l.size
}

func (l *LRUCache[K, V]) PrintList() {
	fmt.Println("[")
	for e := l.list.Front(); e != nil; e = e.Next() {
		fmt.Println(e, e.Value)
	}
	fmt.Println("]")
}

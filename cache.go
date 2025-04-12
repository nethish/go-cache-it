package cache

import (
	"errors"
	"sync"
	"time"
)

var ErrNoEntry = errors.New("ErrNoEntry")

type Value[T any] struct {
	expiration time.Time
	v          T
}

func (v *Value[T]) Expired() bool {
	if v.expiration.IsZero() {
		return false
	}
	return time.Now().After(v.expiration)
}

type SingleCache[K comparable, V any] struct {
	cache map[K]Value[V]
	mu    sync.RWMutex
}

func NewSingleCache[K comparable, V any]() *SingleCache[K, V] {
	return &SingleCache[K, V]{
		cache: make(map[K]Value[V]),
		mu:    sync.RWMutex{},
	}
}

func (c *SingleCache[K, V]) Put(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val := Value[V]{
		v:          v,
		expiration: time.Time{}, // Zero means no expiry
	}
	c.cache[k] = val
}

func (c *SingleCache[K, V]) PutWithExp(k K, v V, dur time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val := Value[V]{
		expiration: time.Now().Add(dur), // Zero means no expiry
		v:          v,
	}
	c.cache[k] = val
}

// Upsert returns true if it writes the value to cache, and false if key already exists
func (c *SingleCache[K, V]) Upsert(k K, v V) bool {
	_, err := c.Get(k)
	if errors.Is(err, ErrNoEntry) {
		c.Put(k, v)
		return true
	}

	return false
}

func (c *SingleCache[K, V]) Get(k K) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.cache[k]
	if !ok {
		return nil, ErrNoEntry
	}

	if val.Expired() {
		return nil, ErrNoEntry
	}

	return val.v, nil
}

func (c *SingleCache[K, V]) Delete(k K) {
	_, err := c.Get(k)
	if errors.Is(err, ErrNoEntry) {
		return
	}

	delete(c.cache, k)
}

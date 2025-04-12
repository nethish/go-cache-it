package cache

import (
	"errors"
	"sync"
	"time"
)

var ErrNoEntry = errors.New("ErrNoEntry")

type Value struct {
	expiration time.Time
	v          any
}

func (v *Value) Expired() bool {
	if v.expiration.IsZero() {
		return false
	}
	return time.Now().After(v.expiration)
}

type SingleCache struct {
	cache map[any]Value
	mu    sync.RWMutex
}

func NewSingleCache() *SingleCache {
	return &SingleCache{
		cache: make(map[any]Value),
		mu:    sync.RWMutex{},
	}
}

func (c *SingleCache) Put(k any, v any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val := Value{
		v:          v,
		expiration: time.Time{}, // Zero means no expiry
	}
	c.cache[k] = val
}

func (c *SingleCache) PutWithExp(k, v any, dur time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val := Value{
		expiration: time.Now().Add(dur), // Zero means no expiry
		v:          v,
	}
	c.cache[k] = val
}

// Upsert returns true if it writes the value to cache, and false if key already exists
func (c *SingleCache) Upsert(k, v any) bool {
	_, err := c.Get(k)
	if errors.Is(err, ErrNoEntry) {
		c.Put(k, v)
		return true
	}

	return false
}

func (c *SingleCache) Get(k any) (any, error) {
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

func (c *SingleCache) Delete(k any) {
	_, err := c.Get(k)
	if errors.Is(err, ErrNoEntry) {
		return
	}

	delete(c.cache, k)
}

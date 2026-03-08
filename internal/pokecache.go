package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entry map[string]cacheEntry
	mu    sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		entry: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return &c
}

func (c *Cache) Add(key string, v []byte) {
	e := cacheEntry{
		createdAt: time.Now(),
		val:       v,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry[key] = e
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.entry[key]
	if !ok {
		return nil, false
	}
	return e.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for key, val := range c.entry {
			if time.Since(val.createdAt) > interval {
				delete(c.entry, key)
			}
		}
		c.mu.Unlock()
	}
}

// Package pokecache manages internal cache logic
package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	// mutex here is important since we will want to read and write to cache concurrently
	// this could be avoided by calling synchronously, but reapLoop will be running constantly. This is a write process
	mu       sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration
}

func (c *Cache) Add(key string, val []byte) {
	// ensure exclusive access to Cache before writing, preventing concurrent writes which could cause a data race or undefined behaviour
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries[key]

	if !ok {
		return nil, false
	} else {
		return entry.val, true
	}
}

func (c *Cache) reapLoop() {
	// all goroutines end when main() closes, or if program is killed or crahes
	// but we still defer closing channel, because if we destroy our cache to open a new one, a running ticker will consume resources
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}
}

func (c *Cache) reap() {
	// write process, so ensure mutex is loked before culling old data
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.entries, key)
		}
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	go c.reapLoop()

	return c
}

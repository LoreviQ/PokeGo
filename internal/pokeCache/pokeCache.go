package pokeCache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.RWMutex{},
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, val []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.cache[key]
	if ok {
		return errors.New("cache already has entry at this key")
	}
	var entry cacheEntry
	entry.val = val
	entry.createdAt = time.Now()
	c.cache[key] = entry
	return nil
}

func (c *Cache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cache[key]
	if !ok {
		return nil, errors.New("cache does not have an entry at this key")
	}
	return entry.val, nil
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		for key, val := range c.cache {
			if time.Since(val.createdAt) >= interval {
				delete(c.cache, key)
				fmt.Printf("Removed - %s\n", key)
			}
		}
		c.mu.Unlock()
	}
}

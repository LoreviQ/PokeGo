package pokeCache

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	cache map[string]cacheEntry
	mu    *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	return Cache{
		cache: make(map[string]cacheEntry),
		mu:    &sync.Mutex{},
	}
}

func (c *Cache) Add(key string, val []byte) error {
	fmt.Println("Adding to cache")
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
	fmt.Println("Reading from cache")
	entry, ok := c.cache[key]
	if !ok {
		return nil, errors.New("cache does not have an entry at this key")
	}
	return entry.val, nil
}

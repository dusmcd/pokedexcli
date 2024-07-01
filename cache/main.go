package cache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntry map[string]cacheEntry
	mu         *sync.Mutex
}

func NewCache(intervalInSeconds int) *Cache {
	// create a new cache
	ticker := time.Tick(time.Duration(intervalInSeconds) * time.Second)
	cache := Cache{
		cacheEntry: make(map[string]cacheEntry),
		mu:         &sync.Mutex{},
	}
	go reapLooper(ticker, &cache, intervalInSeconds)
	return &cache
}

func reapLooper(ticker <-chan time.Time, cache *Cache, intervalInSeconds int) {
	for {
		select {
		case <-ticker:
			// delete cache entries as necessary
			for entry := range cache.cacheEntry {
				cache.mu.Lock()
				timeCreated := cache.cacheEntry[entry].createdAt
				differenceInSeconds := time.Since(timeCreated).Seconds()
				if differenceInSeconds > float64(intervalInSeconds) {
					delete(cache.cacheEntry, entry)
				}
				cache.mu.Unlock()
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (c Cache) GetEntry(key string, ch chan CacheData) {
	// get an entry from cacheEntry map
	c.mu.Lock()
	entry, found := c.cacheEntry[key]
	defer c.mu.Unlock()
	ch <- CacheData{
		Val:   entry.Val,
		Found: found,
	}
}

func (c *Cache) AddEntry(key string, data []byte) {
	c.mu.Lock()
	c.cacheEntry[key] = cacheEntry{
		Val:       data,
		createdAt: time.Now(),
	}
	defer c.mu.Unlock()
}

type cacheEntry struct {
	Val       []byte
	createdAt time.Time
}

type CacheData struct {
	Val   []byte
	Found bool
}

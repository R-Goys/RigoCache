package Rigo

import (
	"github.com/R-Goys/RigoCache/LRU"
	"sync"
)

type cache struct {
	lock       sync.RWMutex
	lru        *LRU.LRUCache
	cacheBytes int64
}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}

func (c *cache) Put(key string, value ByteView) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.lru == nil {
		c.lru = LRU.New(c.cacheBytes, nil)
	}
	c.lru.Put(key, value)
}

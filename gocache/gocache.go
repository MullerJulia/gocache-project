package gocache

import (
	"container/list"
	"sync"
)

type GoCache struct {
	capacity int
	mutex    sync.Mutex
	cache    map[string]*list.Element
	lru      *list.List
}

type entry struct {
	key   string
	value string
}

// NewGoCache creates a new GoCache with the specified capacity.
func NewGoCache(capacity int) *GoCache {
	return &GoCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		lru:      list.New(),
	}
}

// Set adds a key-value pair to the cache or updates the existing value.
func (c *GoCache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, found := c.cache[key]; found {
		c.lru.MoveToFront(elem)
		elem.Value.(*entry).value = value
		return
	}

	if c.lru.Len() >= c.capacity {
		c.evict()
	}

	elem := c.lru.PushFront(&entry{key, value})
	c.cache[key] = elem
}

// Get retrieves the value for a key from the cache.
func (c *GoCache) Get(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, found := c.cache[key]; found {
		c.lru.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}

	return "", false
}

// Delete removes a key-value pair from the cache.
func (c *GoCache) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, found := c.cache[key]; found {
		c.lru.Remove(elem)
		delete(c.cache, key)
		return true
	}

	return false
}

// evict removes the least recently used item from the cache.
func (c *GoCache) evict() {
	elem := c.lru.Back()
	if elem != nil {
		c.lru.Remove(elem)
		delete(c.cache, elem.Value.(*entry).key)
	}
}

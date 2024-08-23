package gocache

import (
	"container/list"
	"sync"
)

type GoCache struct {
	capacity int
	mutex    sync.RWMutex // Use RWMutex to allow multiple readers
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

// Set adds or updates a key-value pair in the cache.
func (c *GoCache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, found := c.cache[key]; found {
		// Update value and move element to front
		elem.Value.(*entry).value = value
		c.lru.MoveToFront(elem)
		return
	}

	// If at capacity, evict the least recently used item
	if c.lru.Len() >= c.capacity {
		c.evict()
	}

	// Add new entry to the front of the list
	elem := c.lru.PushFront(&entry{key, value})
	c.cache[key] = elem
}

// Get retrieves a value from the cache and moves the accessed entry to the front of the LRU list.
func (c *GoCache) Get(key string) (string, bool) {
	// Acquire read lock
	c.mutex.RLock()
	elem, found := c.cache[key]
	if !found {
		// Release read lock early if key is not found
		c.mutex.RUnlock()
		return "", false
	}

	// We found the key, but we need to move it to the front of the LRU list
	value := elem.Value.(*entry).value

	c.mutex.RUnlock() // Release read lock
	c.mutex.Lock()    // Acquire write lock

	// Before moving to the front, check if the element is still present
	// This is a rare case where the element might have been removed in between locks
	if elem, stillFound := c.cache[key]; stillFound {
		c.lru.MoveToFront(elem)
	}

	// Downgrade to read lock for returning the value
	c.mutex.Unlock() // Release write lock
	c.mutex.RLock()  // Acquire read lock again

	defer c.mutex.RUnlock()
	return value, true
}

// evict removes the least recently used item from the cache.
func (c *GoCache) evict() {
	elem := c.lru.Back()
	if elem != nil {
		// Remove from LRU list and cache map
		c.lru.Remove(elem)
		delete(c.cache, elem.Value.(*entry).key)
	}
}

package gocache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGoCache_SetAndGet(t *testing.T) {
	cache := NewGoCache(2)

	cache.Set("a", "apple")
	cache.Set("b", "banana")

	if value, found := cache.Get("a"); !found || value != "apple" {
		t.Errorf("expected apple but got %s", value)
	}

	if value, found := cache.Get("b"); !found || value != "banana" {
		t.Errorf("expected banana but got %s", value)
	}
}

func TestGoCache_Eviction(t *testing.T) {
	cache := NewGoCache(2)

	cache.Set("a", "apple")
	cache.Set("b", "banana")
	cache.Set("c", "cherry") // Should evict "a"

	// Check if 'a' was evicted
	if _, found := cache.Get("a"); found {
		t.Errorf("expected a to be evicted")
	}

	// Ensure 'b' and 'c' are still in the cache
	if value, found := cache.Get("b"); !found || value != "banana" {
		t.Errorf("expected banana but got %s", value)
	}

	if value, found := cache.Get("c"); !found || value != "cherry" {
		t.Errorf("expected cherry but got %s", value)
	}

	// Edge case: attempt to Get an evicted key
	cache.Set("d", "date") // Should evict "b"
	if _, found := cache.Get("b"); found {
		t.Errorf("expected b to be evicted after inserting d")
	}
}

func TestGoCache_Delete(t *testing.T) {
	cache := NewGoCache(2)

	cache.Set("a", "apple")
	cache.Set("b", "banana")

	if deleted := cache.Delete("a"); !deleted {
		t.Errorf("expected a to be deleted")
	}

	if _, found := cache.Get("a"); found {
		t.Errorf("expected a to be not found")
	}
}

func TestGoCache_Concurrency(t *testing.T) {
	cache := NewGoCache(10)
	var wg sync.WaitGroup
	keys := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := keys[i%len(keys)]
			cache.Set(key, key)
			cache.Get(key)
		}(i)
	}

	wg.Wait()

	// Additional verification after concurrent operations
	for _, key := range keys {
		if value, found := cache.Get(key); !found || value != key {
			t.Errorf("expected %s to be found with value %s, but got %v", key, key, value)
		}
	}
}

func TestGoCache_ConcurrentMixedOperations(t *testing.T) {
	cache := NewGoCache(5)
	var wg sync.WaitGroup
	keys := []string{"a", "b", "c", "d", "e"}

	for i := 0; i < 100; i++ {
		wg.Add(3)

		// Concurrent Set
		go func(i int) {
			defer wg.Done()
			key := keys[i%len(keys)]
			cache.Set(key, key)
		}(i)

		// Concurrent Get
		go func(i int) {
			defer wg.Done()
			key := keys[i%len(keys)]
			cache.Get(key)
		}(i)

		// Concurrent Delete
		go func(i int) {
			defer wg.Done()
			key := keys[i%len(keys)]
			cache.Delete(key)
		}(i)
	}

	wg.Wait()

	// Verify final state after mixed operations
	for _, key := range keys {
		value, found := cache.Get(key)
		if found {
			t.Logf("Key: %s, Value: %s\n", key, value)
		} else {
			t.Logf("Key: %s was deleted\n", key)
		}
	}
}

func TestGoCache_RapidUpdates(t *testing.T) {
	cache := NewGoCache(1)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Set("a", "value")
			// Adding a small delay to make sure Set completes
			time.Sleep(10 * time.Millisecond)
			cache.Set("a", "updated-value")
			if value, found := cache.Get("a"); !found || value != "updated-value" {
				t.Errorf("expected updated-value but got %s", value)
			}
		}(i)
	}

	wg.Wait()
}

func TestGoCache_NonExistentKeys(t *testing.T) {
	cache := NewGoCache(2)

	if _, found := cache.Get("nonexistent"); found {
		t.Errorf("expected nonexistent key to return false")
	}

	if deleted := cache.Delete("nonexistent"); deleted {
		t.Errorf("expected delete on nonexistent key to return false")
	}
}

func TestGoCache_StressLRU(t *testing.T) {
	cache := NewGoCache(3)
	var wg sync.WaitGroup
	keys := []string{"a", "b", "c", "d", "e"}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := keys[i%len(keys)]
			cache.Set(key, key)
			cache.Get(key)
		}(i)
	}

	wg.Wait()

	// Perform intermediate checks
	if _, found := cache.Get("d"); !found {
		t.Logf("Key d was not found, possibly evicted as per LRU policy")
	}

	if _, found := cache.Get("e"); !found {
		t.Logf("Key e was not found, possibly evicted as per LRU policy")
	}

	// Final check: log the cache state
	t.Log("Final cache state:")
	for _, key := range keys {
		if value, found := cache.Get(key); found {
			t.Logf("Key: %s, Value: %s\n", key, value)
		} else {
			t.Logf("Key: %s not found, possibly evicted as per LRU policy\n", key)
		}
	}

	// Ensure the cache only contains the correct number of items
	if len(cache.cache) > cache.capacity {
		t.Errorf("Cache exceeds its capacity of %d", cache.capacity)
	}

	// Verify that there is no data corruption or race conditions
	if err := checkCacheIntegrity(cache, keys); err != nil {
		t.Errorf("Cache integrity check failed: %v", err)
	}

	// Edge case: Verify eviction after high concurrency
	cache.Set("f", "fig") // This should evict one of the existing keys
	if _, found := cache.Get("a"); found {
		t.Logf("Key a was found, possibly due to recent access pattern")
	} else {
		t.Logf("Key a was evicted, as expected")
	}
}

// checkCacheIntegrity is a helper function that verifies the cache's internal state
func checkCacheIntegrity(cache *GoCache, keys []string) error {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	// Verify that each element in the LRU list matches the cache map
	for elem := cache.lru.Front(); elem != nil; elem = elem.Next() {
		entry, ok := elem.Value.(*entry)
		if !ok || cache.cache[entry.key] != elem {
			return fmt.Errorf("inconsistent LRU list and cache map for key %s", entry.key)
		}
	}

	// Verify that the cache contains only expected keys
	for key := range cache.cache {
		found := false
		for _, k := range keys {
			if key == k {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unexpected key %s found in cache", key)
		}
	}

	return nil
}

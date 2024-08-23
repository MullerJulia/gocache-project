package gocache

import (
	"sync"
	"testing"
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

	if _, found := cache.Get("a"); found {
		t.Errorf("expected a to be evicted")
	}

	if value, found := cache.Get("b"); !found || value != "banana" {
		t.Errorf("expected banana but got %s", value)
	}

	if value, found := cache.Get("c"); !found || value != "cherry" {
		t.Errorf("expected cherry but got %s", value)
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

	for _, key := range keys {
		if _, found := cache.Get(key); !found {
			t.Errorf("expected %s to be found", key)
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
}

func TestGoCache_RapidUpdates(t *testing.T) {
	cache := NewGoCache(1)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Set("a", "value")
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

	// Verify that the most recently used keys are still present
	if _, found := cache.Get("d"); !found {
		t.Errorf("expected d to be found")
	}
	if _, found := cache.Get("e"); !found {
		t.Errorf("expected e to be found")
	}
}

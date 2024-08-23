package main

import (
	"fmt"

	"github.com/MullerJulia/gocache-project/gocache"
)

func main() {
	// Create a new GoCache with a capacity of 3
	cache := gocache.NewGoCache(3)

	// Set some key-value pairs
	cache.Set("a", "apple")
	cache.Set("b", "banana")
	cache.Set("c", "cherry")

	// Get and print a value
	if value, found := cache.Get("a"); found {
		fmt.Println("Key a:", value)
	}

	// Set another key-value pair, this should trigger an eviction
	cache.Set("d", "date")

	// Check if the least recently used item was evicted
	if _, found := cache.Get("b"); !found {
		fmt.Println("Key b was evicted")
	}

	// Delete a key
	cache.Delete("c")

	// Try to get the deleted key
	if _, found := cache.Get("c"); !found {
		fmt.Println("Key c was deleted")
	}
}
